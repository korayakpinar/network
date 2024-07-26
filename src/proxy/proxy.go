package proxy

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

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
	http.HandleFunc("/", p.proxyHandler)
	log.Printf("Proxy server is running on port %s", p.Port)
	log.Fatal(http.ListenAndServe(":"+p.Port, nil))
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
		var req TransactionRequest

		// Change this response to include the transaction hash after transaction has been decrypted
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"jsonrpc":"2.0","result":"Catched!","id":1}`))

		if err := json.Unmarshal(body, &req); err != nil {
			log.Fatalf("Failed to unmarshal JSON: %v", err)
		}
		fmt.Println("Transaction request taken from the wallet: ", req.Params[0])

		err = p.Handler.HandleTransaction(req.Params[0])
		if err != nil {
			log.Fatalf("Failed to handle transaction: %v", err)
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
