package proxy

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/korayakpinar/network/src/handler"
	"github.com/korayakpinar/network/src/types"
	"github.com/rs/cors"
)

type Proxy struct {
	Handler *handler.Handler
	RpcURL  string
	Port    string
	clients map[*websocket.Conn]bool
	mutex   sync.Mutex
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins
	},
}

type TransactionRequest struct {
	ID      int64    `json:"id"`
	JSONRPC string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
}

type TransactionResponse struct {
	Hash                   string               `json:"hash"`
	Status                 string               `json:"status"`
	RawTx                  types.RawTransaction `json:"rawTx,omitempty"`
	ReceivedAt             time.Time            `json:"receivedAt"`
	ProposedAt             time.Time            `json:"proposedAt,omitempty"`
	DecryptedAt            time.Time            `json:"decryptedAt,omitempty"`
	IncludedAt             time.Time            `json:"includedAt,omitempty"`
	CommitteeSize          int                  `json:"committeeSize"`
	Threshold              int                  `json:"threshold"`
	PartialDecryptionCount int                  `json:"partialDecryptionCount"`
	PartialDecryptions     map[uint64][]byte    `json:"partialDecryptions,omitempty"`
}

func NewProxy(handler *handler.Handler, rpcURL, port string) *Proxy {
	return &Proxy{
		Handler: handler,
		RpcURL:  rpcURL,
		Port:    port,
		clients: make(map[*websocket.Conn]bool),
	}
}

func (p *Proxy) Start() {
	r := mux.NewRouter()
	r.HandleFunc("/ws", p.handleWebSocket)
	r.HandleFunc("/", p.proxyHandler)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	log.Printf("Proxy server is running on port %s", p.Port)
	log.Fatal(http.ListenAndServe(":"+p.Port, handler))
}

func (p *Proxy) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	p.mutex.Lock()
	p.clients[conn] = true
	p.mutex.Unlock()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		var req map[string]interface{}
		if err := json.Unmarshal(message, &req); err != nil {
			log.Println(err)
			continue
		}

		switch req["type"] {
		case "getTxStatus":
			p.handleTxStatus(conn, req["txHash"].(string))
		case "getRecentTransactions":
			p.handleRecentTransactions(conn)
		case "subscribe":
			go p.handleSubscription(conn, req["txHash"].(string))
		}

		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Println(err)
			break
		}
	}

	p.mutex.Lock()
	delete(p.clients, conn)
	p.mutex.Unlock()
}

func (p *Proxy) handleTxStatus(conn *websocket.Conn, txHash string) {
	tx := p.Handler.GetTransaction(txHash)
	if tx == nil {
		conn.WriteJSON(map[string]string{"error": "Transaction not found"})
		return
	}

	response := p.createTransactionResponse(tx)
	conn.WriteJSON(response)
}

func (p *Proxy) handleRecentTransactions(conn *websocket.Conn) {
	transactions := p.Handler.GetRecentTransactions()

	var responses []TransactionResponse
	for _, tx := range transactions {
		responses = append(responses, p.createTransactionResponse(tx))
	}

	conn.WriteJSON(map[string]interface{}{
		"type":         "recentTransactions",
		"transactions": responses,
	})
}

func (p *Proxy) handleSubscription(conn *websocket.Conn, txHash string) {
	for {
		tx := p.Handler.GetTransaction(txHash)
		if tx == nil {
			time.Sleep(time.Second)
			continue
		}

		response := p.createTransactionResponse(tx)
		if err := conn.WriteJSON(response); err != nil {
			log.Println(err)
			return
		}

		if tx.Status == types.StatusIncluded {
			return
		}

		time.Sleep(time.Second)
	}
}

func (p *Proxy) createTransactionResponse(tx *types.Transaction) TransactionResponse {
	status := "pending"
	switch tx.Status {
	case types.StatusProposed:
		status = "proposed"
	case types.StatusDecrypted:
		status = "decrypted"
	case types.StatusIncluded:
		status = "included"
	}

	partDecs := p.Handler.GetPartialDecryptions(tx.Hash)

	if tx.Status == types.StatusDecrypted || tx.Status == types.StatusIncluded && tx.RawTransaction != nil {
		rawTx, err := types.DecodeRawTransaction(tx.DecryptedTransaction.Body.Content[2:])
		if err != nil {
			log.Printf("Failed to decode raw transaction: %v", err)
		}
		return TransactionResponse{
			Hash:                   tx.Hash,
			Status:                 status,
			RawTx:                  *rawTx,
			ReceivedAt:             tx.ReceivedAt,
			ProposedAt:             tx.ProposedAt,
			DecryptedAt:            tx.DecryptedAt,
			IncludedAt:             tx.IncludedAt,
			CommitteeSize:          tx.CommitteeSize + 1,
			Threshold:              tx.Threshold,
			PartialDecryptionCount: tx.PartialDecryptionCount,
			PartialDecryptions:     partDecs,
		}
	} else {
		return TransactionResponse{
			Hash:                   tx.Hash,
			Status:                 status,
			ReceivedAt:             tx.ReceivedAt,
			ProposedAt:             tx.ProposedAt,
			DecryptedAt:            tx.DecryptedAt,
			IncludedAt:             tx.IncludedAt,
			CommitteeSize:          tx.CommitteeSize + 1,
			Threshold:              tx.Threshold,
			PartialDecryptionCount: tx.PartialDecryptionCount,
			PartialDecryptions:     partDecs,
		}
	}
}

func (p *Proxy) proxyHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	var req map[string]interface{}
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	method, ok := req["method"].(string)
	if !ok {
		http.Error(w, "Method not found in request", http.StatusBadRequest)
		return
	}

	if method == "eth_sendTransaction" || method == "eth_sendRawTransaction" {
		var txReq TransactionRequest
		if err := json.Unmarshal(body, &txReq); err != nil {
			log.Printf("Failed to unmarshal JSON: %v", err)
			http.Error(w, "Invalid JSON request", http.StatusBadRequest)
			return
		}

		rawTx := txReq.Params[0]
		txHash, err := p.Handler.CalculateTxHash(rawTx)
		if err != nil {
			http.Error(w, "Failed to calculate transaction hash", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"jsonrpc": "2.0",
			"result":  txHash,
			"id":      txReq.ID,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

		if !p.Handler.CheckTransactionDuplicate(txHash) {
			err = p.Handler.HandleTransaction(rawTx)
			if err != nil {
				log.Printf("Failed to handle transaction: %v", err)
			}
		}

		return
	}

	respBody, statusCode, err := p.forwardRequest(body)
	if err != nil {
		http.Error(w, "Failed to forward request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(respBody)
}

func (p *Proxy) forwardRequest(body []byte) ([]byte, int, error) {
	req, err := http.NewRequest("POST", p.RpcURL, strings.NewReader(string(body)))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return responseBody, resp.StatusCode, nil
}
