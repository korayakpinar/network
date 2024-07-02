package api

import (
	"bytes"
	"net/http"

	"google.golang.org/protobuf/proto"
)

func DecryptData(enc string, pks []string, sa1 string, sa2 string, iv string, t uint64, n uint64) (string, error) {
	client := http.Client{}
	req := &DecryptDataRequest{
		Enc: enc,
		Pks: pks,
		Sa1: sa1,
		Sa2: sa2,
		Iv:  iv,
		T:   t,
		N:   n,
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

	var decryptDataResp DecryptDataResponse
	err = proto.Unmarshal(bodyBytes, &decryptDataResp)
	if err != nil {
		return "", err
	}

	return decryptDataResp.Result, nil
}

func PartialDecrypt(gammaG2 string) (string, error) {
	client := http.Client{}

	req := &PartialDecryptRequest{
		GammaG2: gammaG2,
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return "", err
	}

	postReader := bytes.NewReader(data)

	resp, err := client.Post("http://127.0.0.1:8080/partdec", "application/x-www-form-urlencoded", postReader)
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

	var partDecResp PartialDecryptResponse
	err = proto.Unmarshal(bodyBytes, &partDecResp)
	if err != nil {
		return "", err
	}

	return partDecResp.Result, nil
}

func VerifyPart(pk string, gammaG2 string, partDec string) (string, error) {
	client := http.Client{}

	req := &VerifyPartRequest{
		Pk:      pk,
		GammaG2: gammaG2,
		PartDec: partDec,
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

	var verifyPartResp VerifyPartResponse
	err = proto.Unmarshal(bodyBytes, &verifyPartResp)
	if err != nil {
		return "", err
	}

	return verifyPartResp.Result, nil
}
