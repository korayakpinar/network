package crypto

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/protobuf/proto"
)

type Crypto struct {
	port string
}

func NewCrypto(port string) *Crypto {
	return &Crypto{port: port}
}

func (c *Crypto) EncryptTransaction(msg []byte, pks [][]byte, t uint64, n uint64) (*EncryptResponse, error) {
	client := http.Client{}

	req := &EncryptRequest{
		Msg: msg,
		Pks: pks,
		T:   t,
		N:   n,
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	postReader := bytes.NewReader(data)

	url := fmt.Sprintf("http://127.0.0.1:%s/encrypt", c.port)
	resp, err := client.Post(url, "application/protobuf", postReader)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var encryptDataResp EncryptResponse
	err = proto.Unmarshal(bodyBytes, &encryptDataResp)
	if err != nil {
		return nil, err
	}

	return &encryptDataResp, nil
}

func (c *Crypto) DecryptTransaction(enc []byte, pks [][]byte, parts map[uint64][]byte, gammaG2 []byte, sa1 []byte, sa2 []byte, iv []byte, t uint64, n uint64) ([]byte, error) {
	client := http.Client{}

	req := &DecryptRequest{
		Enc:     []byte(enc),
		Pks:     pks,
		Parts:   parts,
		GammaG2: gammaG2,
		Sa1:     sa1,
		Sa2:     sa2,
		Iv:      iv,
		T:       t,
		N:       n,
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	postReader := bytes.NewReader(data)

	url := fmt.Sprintf("http://127.0.0.1:%s/decrypt", c.port)
	resp, err := client.Post(url, "application/protobuf", postReader)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var decryptDataResp Response
	err = proto.Unmarshal(bodyBytes, &decryptDataResp)
	if err != nil {
		return nil, err
	}

	return decryptDataResp.Result, nil
}

func (c *Crypto) PartialDecrypt(gammaG2 []byte) ([]byte, error) {
	client := http.Client{}

	req := &GammaG2Request{
		GammaG2: []byte(gammaG2),
	}
	data, err := proto.Marshal(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	postReader := bytes.NewReader(data)

	url := fmt.Sprintf("http://127.0.0.1:%s/partdec", c.port)
	resp, err := client.Post(url, "application/protobuf", postReader)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var partDecResp Response
	err = proto.Unmarshal(bodyBytes, &partDecResp)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return partDecResp.Result, nil
}

func (c *Crypto) GetPK(id uint64, n uint64) ([]byte, error) {
	client := http.Client{}

	req := &PKRequest{
		Id: id,
		N:  n,
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	postReader := bytes.NewReader(data)

	url := fmt.Sprintf("http://127.0.0.1:%s/getpk", c.port)
	resp, err := client.Post(url, "application/protobuf", postReader)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var partDecResp Response
	err = proto.Unmarshal(bodyBytes, &partDecResp)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return partDecResp.Result, nil
}

func (c *Crypto) VerifyPart(pk []byte, gammaG2 []byte, partDec []byte) error {
	client := http.Client{}

	req := &VerifyPartRequest{
		Pk:      []byte(pk),
		GammaG2: []byte(gammaG2),
		PartDec: []byte(partDec),
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return err
	}

	postReader := bytes.NewReader(data)

	url := fmt.Sprintf("http://127.0.0.1:%s/verifypart", c.port)
	resp, err := client.Post(url, "application/protobuf", postReader)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return err
	}
	resp.Body.Close()

	return nil
}
