package ipfs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
)

const (
	apiPort     = 5001
	gatewayPort = 8080
)

type IPFSService struct {
	APIBaseURL string
	GatewayURL string
}

// NewIPFSService creates a new IPFSService with hardcoded ports
func NewIPFSService(baseURL string) *IPFSService {
	// Ensure baseURL doesn't end with a slash
	baseURL = strings.TrimRight(baseURL, "/")

	return &IPFSService{
		APIBaseURL: fmt.Sprintf("%s:%d", baseURL, apiPort),
		GatewayURL: fmt.Sprintf("%s:%d", baseURL, gatewayPort),
	}
}

// CalculateCID calculates the CID for the given content
func (is *IPFSService) CalculateCID(content []byte) (string, error) {
	pref := cid.Prefix{
		Version:  0,
		Codec:    cid.DagProtobuf,
		MhType:   multihash.SHA2_256,
		MhLength: -1, // Default length
	}

	c, err := pref.Sum(content)
	if err != nil {
		return "", fmt.Errorf("error creating CID: %w", err)
	}

	return c.String(), nil
}

// UploadKey uploads a key as bytes to IPFS
func (is *IPFSService) UploadKey(key []byte) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "key")
	if err != nil {
		return "", fmt.Errorf("error creating form file: %w", err)
	}

	_, err = part.Write(key)
	if err != nil {
		return "", fmt.Errorf("error writing key content: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("error closing multipart writer: %w", err)
	}

	req, err := http.NewRequest("POST", is.APIBaseURL+"/api/v0/add?cid-version=0", body)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result struct {
		Hash string `json:"Hash"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	return result.Hash, nil
}

// GetKeyByCID retrieves a key from IPFS using its CID
func (is *IPFSService) GetKeyByCID(cid string) ([]byte, error) {
	url := fmt.Sprintf("%s/ipfs/%s", is.GatewayURL, cid)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
