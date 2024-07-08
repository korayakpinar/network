package crypto

import (
	"bytes"
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
		return []byte{}, err
	}

	postReader := bytes.NewReader(data)

	resp, err := client.Post("http://127.0.0.1:8080/encrypt", "application/x-www-form-urlencoded", postReader)

	if err != nil {
		return []byte{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return []byte{}, err
	}
	defer resp.Body.Close()

	var bodyBytes []byte
	_, err = resp.Body.Read(bodyBytes)
	if err != nil {
		return []byte{}, err
	}

	var encryptDataResp Response
	err = proto.Unmarshal(bodyBytes, &encryptDataResp)
	if err != nil {
		return []byte{}, err
	}

	return encryptDataResp.Result, nil
}
func DecryptTransaction(enc []byte, pks [][]byte, parts [][]byte, sa1 []byte, sa2 []byte, iv []byte, t uint64, n uint64) (string, error) {
	client := http.Client{}

	req := &DecryptDataRequest{
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
		return "", err
	}

	postReader := bytes.NewReader(data)

	resp, err := client.Post("http://127.0.0.1:8080/decrypt", "application/x-www-form-urlencoded", postReader)

	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", err
	}
	defer resp.Body.Close()

	var bodyBytes []byte
	_, err = resp.Body.Read(bodyBytes)
	if err != nil {
		return "", err
	}

	var decryptDataResp Response
	err = proto.Unmarshal(bodyBytes, &decryptDataResp)
	if err != nil {
		return "", err
	}

	return string(decryptDataResp.Result), nil
}

func PartialDecrypt(gammaG2 []byte) ([]byte, error) {
	client := http.Client{}

	req := &PartialDecryptRequest{
		GammaG2: []byte(gammaG2),
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return []byte{}, err
	}

	postReader := bytes.NewReader(data)

	resp, err := client.Post("http://127.0.0.1:8080/partdec", "application/x-www-form-urlencoded", postReader)
	if err != nil {
		return []byte{}, err
	}
	if resp.StatusCode != 200 {
		return []byte{}, err
	}
	defer resp.Body.Close()

	var bodyBytes []byte
	_, err = resp.Body.Read(bodyBytes)
	if err != nil {
		return []byte{}, err
	}

	var partDecResp Response
	err = proto.Unmarshal(bodyBytes, &partDecResp)
	if err != nil {
		return []byte{}, err
	}

	return partDecResp.Result, nil
}

func VerifyPart(pk []byte, gammaG2 []byte, partDec []byte) (string, error) {
	client := http.Client{}

	req := &VerifyPartRequest{
		Pk:      []byte(pk),
		GammaG2: []byte(gammaG2),
		PartDec: []byte(partDec),
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return "", err
	}

	postReader := bytes.NewReader(data)

	resp, err := client.Post("http://127.0.0.1:8080/verifypart", "application/x-www-form-urlencoded", postReader)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", err
	}
	defer resp.Body.Close()

	var bodyBytes []byte
	_, err = resp.Body.Read(bodyBytes)
	if err != nil {
		return "", err
	}

	var verifyPartResp Response
	err = proto.Unmarshal(bodyBytes, &verifyPartResp)
	if err != nil {
		return "", err
	}

	return string(verifyPartResp.Result), nil
}
