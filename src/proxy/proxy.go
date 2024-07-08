package proxy

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// Proxy represents the proxy server.
type Proxy struct {
	RpcURL string
	Port   string
}

// NewProxy creates a new Proxy instance.
func NewProxy(rpcURL, port string) *Proxy {
	return &Proxy{
		RpcURL: rpcURL,
		Port:   port,
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
		fmt.Println("Catched!")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"jsonrpc":"2.0","result":"Catched!","id":1}`))

		// --- TODO ---
		// Encrypt the transaction data with the encrypt function
		// Call the Client's add method to add the encrypted transaction to the mempool
		// and broadcast it to the network

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
