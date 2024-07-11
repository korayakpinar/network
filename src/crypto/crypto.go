package crypto

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/protobuf/proto"
)

func EncryptTransaction(msg []byte, pks [][]byte, t uint64, n uint64) ([]byte, error) {
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

	resp, err := client.Post("http://127.0.0.1:8080/encrypt", "application/protobuf", postReader)

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

	var encryptDataResp Response
	err = proto.Unmarshal(bodyBytes, &encryptDataResp)
	if err != nil {
		return nil, err
	}

	return encryptDataResp.Result, nil
}

func DecryptTransaction(enc []byte, pks [][]byte, parts [][]byte, sa1 []byte, sa2 []byte, iv []byte, t uint64, n uint64) ([]byte, error) {
	client := http.Client{}

	req := &DecryptParamsRequest{
		Enc:   []byte(enc),
		Pks:   pks,
		Parts: parts,
		Sa1:   sa1,
		Sa2:   sa2,
		Iv:    iv,
		T:     t,
		N:     n,
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	postReader := bytes.NewReader(data)

	resp, err := client.Post("http://127.0.0.1:8080/decrypt", "application/protobuf", postReader)

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

func PartialDecrypt(gammaG2 []byte) ([]byte, error) {
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

	resp, err := client.Post("http://127.0.0.1:8080/partdec", "application/protobuf", postReader)
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

func VerifyPart(pk []byte, gammaG2 []byte, partDec []byte) error {
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

	resp, err := client.Post("http://127.0.0.1:8080/verifydec", "application/protobuf", postReader)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return err
	}
	resp.Body.Close()

	return nil
}
