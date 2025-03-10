package proxy

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/korayakpinar/network/src/mempool"
	"github.com/korayakpinar/network/src/types"
	"github.com/korayakpinar/network/src/utils"
	"github.com/rs/cors"
)

type Proxy struct {
	Mempool                         *mempool.Mempool
	Handler                         types.Handler
	RpcURL                          string
	Port                            string
	clients                         map[*websocket.Conn]bool
	subscriptions                   map[*websocket.Conn]map[string]bool
	recentTransactionsSubscriptions map[*websocket.Conn]bool
	mutex                           sync.Mutex
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

type EncryptedTransactionRequest struct {
	ID      int64  `json:"id"`
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []struct {
		Hash        string   `json:"hash"`
		EncryptedTx []byte   `json:"encryptedTx"`
		PkIDs       []uint64 `json:"pkIDs"`
		GammaG2     []byte   `json:"gammaG2"`
		Threshold   uint64   `json:"threshold"`
		Sa1         []byte   `json:"sa1"`
		Sa2         []byte   `json:"sa2"`
		Iv          []byte   `json:"iv"`
	} `json:"params"`
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
	AlreadyEncrypted       bool                 `json:"alreadyEncrypted"`
}

type RecentTransactionsUpdate struct {
	Type         string                `json:"type"`
	Transactions []TransactionResponse `json:"transactions"`
}

func NewProxy(handler types.Handler, mempool *mempool.Mempool, rpcURL, port string) *Proxy {
	return &Proxy{
		Handler:                         handler,
		Mempool:                         mempool,
		RpcURL:                          rpcURL,
		Port:                            port,
		clients:                         make(map[*websocket.Conn]bool),
		subscriptions:                   make(map[*websocket.Conn]map[string]bool),
		recentTransactionsSubscriptions: make(map[*websocket.Conn]bool),
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

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	server := &http.Server{
		Addr:      ":" + p.Port,
		Handler:   handler,
		TLSConfig: tlsConfig,
	}

	log.Printf("HTTPS proxy server is running on port %s", p.Port)
	log.Fatal(server.ListenAndServeTLS("/app/cert.pem", "/app/key.pem"))
}

func (p *Proxy) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	p.mutex.Lock()
	p.clients[conn] = true
	p.subscriptions[conn] = make(map[string]bool)
	p.mutex.Unlock()

	defer func() {
		p.mutex.Lock()
		delete(p.clients, conn)
		delete(p.subscriptions, conn)
		delete(p.recentTransactionsSubscriptions, conn)
		p.mutex.Unlock()
		conn.Close()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		var req map[string]interface{}
		if err := json.Unmarshal(message, &req); err != nil {
			log.Printf("JSON unmarshal error: %v", err)
			continue
		}

		switch req["type"] {
		case "getTxStatus":
			p.handleTxStatus(conn, req["txHash"].(string))
		case "getRecentTransactions":
			p.handleRecentTransactions(conn)
		case "subscribe":
			p.handleSubscription(conn, req["txHash"].(string))
		case "subscribeToNewTransactions":
			p.handleNewTransactionsSubscription(conn)
		}
	}

	log.Println("WebSocket connection closed")
}

func (p *Proxy) handleNewTransactionsSubscription(conn *websocket.Conn) {
	p.mutex.Lock()
	p.subscriptions[conn]["newTransactions"] = true
	p.mutex.Unlock()
}

func (p *Proxy) BroadcastNewTransaction(tx *types.Transaction) {
	response := p.createTransactionResponse(tx)
	message := map[string]interface{}{
		"type":        "newTransaction",
		"transaction": response,
	}

	p.mutex.Lock()
	for conn, subs := range p.subscriptions {
		if subs["newTransactions"] {
			err := conn.WriteJSON(message)
			if err != nil {
				log.Printf("Error sending new transaction to client: %v", err)
				delete(p.clients, conn)
				delete(p.subscriptions, conn)
				conn.Close()
			}
		}
	}
	p.mutex.Unlock()

	p.BroadcastRecentTransactionsUpdate()
}

func (p *Proxy) handleTxStatus(conn *websocket.Conn, txHash string) {
	tx := p.Mempool.GetTransaction(txHash)
	if tx == nil {
		conn.WriteJSON(map[string]string{"error": "Transaction not found"})
		return
	}

	response := p.createTransactionResponse(tx)
	conn.WriteJSON(response)
}

func (p *Proxy) handleRecentTransactions(conn *websocket.Conn) {
	transactions := p.Mempool.GetRecentTransactions()

	var responses []TransactionResponse
	for _, tx := range transactions {
		responses = append(responses, p.createTransactionResponse(tx))
	}

	update := RecentTransactionsUpdate{
		Type:         "recentTransactions",
		Transactions: responses,
	}

	conn.WriteJSON(update)

	// Subscribe to real-time updates
	p.mutex.Lock()
	p.recentTransactionsSubscriptions[conn] = true
	p.mutex.Unlock()
}

func (p *Proxy) handleSubscription(conn *websocket.Conn, txHash string) {
	p.mutex.Lock()
	p.subscriptions[conn][txHash] = true
	p.mutex.Unlock()

	go p.streamTransactionUpdates(conn, txHash)
}

func (p *Proxy) streamTransactionUpdates(conn *websocket.Conn, txHash string) {
	ticker := time.NewTicker(time.Second / 4)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.mutex.Lock()
			if !p.subscriptions[conn][txHash] {
				p.mutex.Unlock()
				return
			}
			p.mutex.Unlock()

			tx := p.Mempool.GetTransaction(txHash)
			if tx == nil {
				continue
			}

			response := p.createTransactionResponse(tx)
			message := map[string]interface{}{
				"type":   "txUpdate",
				"txHash": txHash,
				"data":   response,
			}

			if err := conn.WriteJSON(message); err != nil {
				log.Printf("Error sending transaction update: %v", err)
				p.mutex.Lock()
				delete(p.subscriptions[conn], txHash)
				p.mutex.Unlock()
				return
			}

		}
	}
}

func (p *Proxy) BroadcastTransactionUpdate(tx *types.Transaction) {
	response := p.createTransactionResponse(tx)
	message := map[string]interface{}{
		"type":   "txUpdate",
		"txHash": tx.Hash,
		"data":   response,
	}

	p.mutex.Lock()
	for conn, subs := range p.subscriptions {
		if subs[tx.Hash] {
			err := conn.WriteJSON(message)
			if err != nil {
				log.Printf("Error sending transaction update to client: %v", err)
				delete(p.clients, conn)
				delete(p.subscriptions, conn)
				conn.Close()
			}
		}
	}
	p.mutex.Unlock()

	p.BroadcastRecentTransactionsUpdate()
}

func (p *Proxy) BroadcastRecentTransactionsUpdate() {
	transactions := p.Mempool.GetRecentTransactions()

	var responses []TransactionResponse
	for _, tx := range transactions {
		responses = append(responses, p.createTransactionResponse(tx))
	}

	update := RecentTransactionsUpdate{
		Type:         "recentTransactionsUpdate",
		Transactions: responses,
	}

	p.mutex.Lock()
	for conn := range p.recentTransactionsSubscriptions {
		err := conn.WriteJSON(update)
		if err != nil {
			log.Printf("Error sending recent transactions update to client: %v", err)
			delete(p.recentTransactionsSubscriptions, conn)
			conn.Close()
		}
	}
	p.mutex.Unlock()
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

	partDecs := p.Mempool.GetPartialDecryptions(tx.Hash)

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
			AlreadyEncrypted:       tx.AlreadyEncrypted,
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
			AlreadyEncrypted:       tx.AlreadyEncrypted,
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

	switch method {
	case "eth_sendTransaction", "eth_sendRawTransaction":
		p.handleRawTransaction(w, body)
	case "eth_sendEncryptedTransaction":
		p.handleEncryptedTransaction(w, body)
	default:
		p.forwardToRPC(w, body)
	}
}

func (p *Proxy) handleRawTransaction(w http.ResponseWriter, body []byte) {
	var txReq TransactionRequest
	if err := json.Unmarshal(body, &txReq); err != nil {
		log.Printf("Failed to unmarshal JSON: %v", err)
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	rawTx := txReq.Params[0]
	txHash, err := utils.CalculateTxHash(rawTx)
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

	if !p.Mempool.CheckTransactionDuplicate(txHash) {
		err = p.Handler.HandleTransaction(rawTx)
		if err != nil {
			log.Printf("Failed to handle transaction: %v", err)
		}
		p.BroadcastNewTransaction(p.Mempool.GetTransaction(txHash))
	}
}

func (p *Proxy) handleEncryptedTransaction(w http.ResponseWriter, body []byte) {
	var txReq EncryptedTransactionRequest
	if err := json.Unmarshal(body, &txReq); err != nil {
		log.Printf("Failed to unmarshal encrypted transaction JSON: %v", err)
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	if len(txReq.Params) == 0 {
		http.Error(w, "No parameters provided", http.StatusBadRequest)
		return
	}

	encTx := txReq.Params[0]

	// Use the hash from the request
	txHash := encTx.Hash

	// Create types.EncryptedTransaction
	encryptedTransaction := &types.EncryptedTransaction{
		Header: &types.EncryptedTxHeader{
			Hash:    txHash,
			GammaG2: encTx.GammaG2,
			PkIDs:   encTx.PkIDs,
		},
		Body: &types.EncryptedTxBody{
			Sa1:       encTx.Sa1,
			Sa2:       encTx.Sa2,
			Iv:        encTx.Iv,
			EncText:   encTx.EncryptedTx,
			Threshold: encTx.Threshold,
		},
	}

	fmt.Println(txHash)
	fmt.Println(encTx.GammaG2)
	fmt.Println(encTx.PkIDs)
	fmt.Println(len(encTx.PkIDs))
	fmt.Println(encTx.Threshold)
	fmt.Println(encTx.EncryptedTx)

	// Prepare the response
	response := map[string]interface{}{
		"jsonrpc": "2.0",
		"result":  txHash,
		"id":      txReq.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	// Handle the encrypted transaction if it's not a duplicate
	if !p.Mempool.CheckTransactionDuplicate(txHash) {
		err := p.Handler.HandleEncryptedTransaction(encryptedTransaction)
		if err != nil {
			log.Printf("Failed to handle encrypted transaction: %v", err)
		}
		p.BroadcastNewTransaction(p.Mempool.GetTransaction(txHash))
	}
}

func (p *Proxy) forwardToRPC(w http.ResponseWriter, body []byte) {
	req, err := http.NewRequest("POST", p.RpcURL, bytes.NewReader(body))
	if err != nil {
		http.Error(w, "Failed to create forward request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to forward request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(responseBody)
}
