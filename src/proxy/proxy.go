package proxy

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/korayakpinar/network/src/handler"
)

// Proxy represents the proxy server.
type Proxy struct {
	Handler *handler.Handler
	RpcURL  string
	Port    string
}

// TransactionRequest represents the JSON-RPC request for sending a transaction which is taken from the wallet.
type TransactionRequest struct {
	ID      int64    `json:"id"`
	JSONRPC string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
}

// TransactionStatus represents the status of a transaction.
type TransactionStatus struct {
	Proposed               bool   `json:"proposed"`
	PartialDecryptionCount int    `json:"partialDecryptionCount"`
	Decrypted              bool   `json:"decrypted"`
	Included               bool   `json:"included"`
	TxInfo                 TxInfo `json:"txInfo"`
}

// TxInfo represents additional information about the transaction.
type TxInfo struct {
	Hash          string `json:"hash"`
	CommitteeSize int    `json:"committeeSize"`
	Threshold     int    `json:"threshold"`
}

// NewProxy creates a new Proxy instance.
func NewProxy(handler *handler.Handler, rpcURL, port string) *Proxy {
	return &Proxy{
		Handler: handler,
		RpcURL:  rpcURL,
		Port:    port,
	}
}

// Start starts the proxy server.
func (p *Proxy) Start() {
	r := mux.NewRouter()
	r.HandleFunc("/", p.proxyHandler)
	r.HandleFunc("/tx/{txHash}", p.txStatusHandler)
	log.Printf("Proxy server is running on port %s", p.Port)
	log.Fatal(http.ListenAndServe(":"+p.Port, r))
}

// txStatusHandler handles requests for transaction status.
func (p *Proxy) txStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	txHash := vars["txHash"]

	if txHash == "" {
		http.Error(w, "Transaction hash is required", http.StatusBadRequest)
		return
	}

	committeeSize, threshold := p.Handler.GetMetadataOfTx(txHash)

	// Create the response with the transaction status
	status := TransactionStatus{
		Proposed:               p.Handler.CheckTransactionProposed(txHash),
		PartialDecryptionCount: p.Handler.GetPartialDecryptionCount(txHash),
		Decrypted:              p.Handler.CheckTransactionDecrypted(txHash),
		Included:               true,
		TxInfo: TxInfo{
			Hash:          txHash,
			CommitteeSize: committeeSize,
			Threshold:     threshold,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// proxyHandler handles incoming HTTP requests and forwards them to the RPC server.
func (p *Proxy) proxyHandler(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	fmt.Println(string(body))

	// Parse the JSON request data
	var req map[string]interface{}
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	// Check the method name
	method, ok := req["method"].(string)
	if !ok {
		http.Error(w, "Method not found in request", http.StatusBadRequest)
		return
	}
	// Check if it's one of the specified methods
	if method == "eth_sendTransaction" || method == "eth_sendRawTransaction" {
		var txReq TransactionRequest

		if err := json.Unmarshal(body, &txReq); err != nil {
			log.Printf("Failed to unmarshal JSON: %v", err)
			http.Error(w, "Invalid JSON request", http.StatusBadRequest)
			return
		}

		fmt.Println("Transaction request taken from the wallet: ", txReq.Params[0])

		// Calculate the transaction hash
		rawTx := txReq.Params[0]
		if !ok {
			http.Error(w, "Invalid transaction data", http.StatusBadRequest)
			return
		}

		txHash, err := p.Handler.CalculateTxHash(rawTx)
		if err != nil {
			http.Error(w, "Failed to calculate transaction hash", http.StatusInternalServerError)
			return
		}

		// Create the response with the calculated hash
		response := map[string]interface{}{
			"jsonrpc": "2.0",
			"result":  txHash,
			"id":      txReq.ID,
		}

		// Send the response back to the client
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

		if !p.Handler.CheckTransactionDuplicate(txReq.Params[0]) {
			// Handle the transaction (you may want to do this asynchronously)
			err = p.Handler.HandleTransaction(txReq.Params[0])
			if err != nil {
				log.Printf("Failed to handle transaction: %v", err)
			}
		}

		return
	}

	// Forward the request to the RPC server if it's not one of the specified methods
	respBody, statusCode, err := p.forwardRequest(body)
	if err != nil {
		http.Error(w, "Failed to forward request", http.StatusInternalServerError)
		return
	}

	// Send the response back to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(respBody)
}

// forwardRequest forwards the request to the RPC server and returns the response.
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
