package pinata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
)

type IPFSService struct {
	BearerToken string
	APIBaseURL  string
	GatewayURL  string
}

type PinataResponse struct {
	IpfsHash    string `json:"IpfsHash"`
	PinSize     int    `json:"PinSize"`
	Timestamp   string `json:"Timestamp"`
	IsDuplicate bool   `json:"isDuplicate"`
}

func NewIPFSService(bearerToken, gatewayURL string) *IPFSService {
	return &IPFSService{
		BearerToken: bearerToken,
		APIBaseURL:  "https://api.pinata.cloud",
		GatewayURL:  gatewayURL,
	}
}

// CalculateCID, belirtilen JSON içeriği için CID'yi hesaplar
func (ps *IPFSService) CalculateCID(content interface{}) (string, error) {
	// JSON'ı özel bir şekilde serialize et
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "")
	if err := encoder.Encode(content); err != nil {
		return "", fmt.Errorf("JSON dönüşüm hatası: %w", err)
	}

	// Sondaki newline karakterini kaldır
	jsonData := bytes.TrimSuffix(buf.Bytes(), []byte("\n"))

	// CIDv1 oluştur
	pref := cid.Prefix{
		Version:  1,
		Codec:    cid.Raw,
		MhType:   multihash.SHA2_256,
		MhLength: -1, // Varsayılan uzunluk
	}

	c, err := pref.Sum(jsonData)
	if err != nil {
		return "", fmt.Errorf("CID oluşturma hatası: %w", err)
	}

	return c.String(), nil
}

func (ps *IPFSService) UploadJSON(content interface{}, name string) (*PinataResponse, error) {
	url := ps.APIBaseURL + "/pinning/pinJSONToIPFS"

	requestBody := map[string]interface{}{
		"pinataOptions": map[string]interface{}{
			"cidVersion": 1,
		},
		"pinataMetadata": map[string]string{
			"name": name,
		},
		"pinataContent": content,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("JSON dönüşüm hatası: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("istek oluşturma hatası: %w", err)
	}

	ps.setHeaders(req)

	return ps.sendRequest(req)
}

func (ps *IPFSService) GetFileByCID(cid string) ([]byte, error) {
	url := fmt.Sprintf("%s/ipfs/%s", ps.GatewayURL, cid)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("istek oluşturma hatası: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("istek gönderme hatası: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dosya alma hatası: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

// setHeaders, isteklere gerekli başlıkları ekler
func (ps *IPFSService) setHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ps.BearerToken)
}

// sendRequest, Pinata API'sine istek gönderir ve yanıtı işler
func (ps *IPFSService) sendRequest(req *http.Request) (*PinataResponse, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("istek gönderme hatası: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("yanıt okuma hatası: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API hatası: %s", string(body))
	}

	var pinataResp PinataResponse
	err = json.Unmarshal(body, &pinataResp)
	if err != nil {
		return nil, fmt.Errorf("yanıt ayrıştırma hatası: %w", err)
	}

	return &pinataResp, nil
}
