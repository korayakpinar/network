package api

import (
	"net/http"
	"strconv"
	"strings"
)

func DecryptData(enc string, pks []string, sa1 string, sa2 string, iv string, t int, n int) (string, error) {
	client := http.Client{}

	body := ""
	for _, pk := range pks {
		body += pk + " "
	}
	body += sa1 + " " + sa2 + " " + iv + " " + strconv.Itoa(t) + " " + strconv.Itoa(n)
	postReader := strings.NewReader(body)
	resp, err := client.Post("http://127.0.0.1:8080/decrypt", "application/x-www-form-urlencoded", postReader)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", err
	}

	var bodyBytes []byte
	_, err = resp.Body.Read(bodyBytes)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

func PartialDecrypt(gammaG2 string) (string, error) {
	client := http.Client{}

	postReader := strings.NewReader(gammaG2)

	resp, err := client.Post("http://127.0.0.1:8080/partdec", "application/x-www-form-urlencoded", postReader)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", err
	}

	var bodyBytes []byte
	_, err = resp.Body.Read(bodyBytes)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil
}

func VerifyPart(pk string, gammaG2 string, partDec string) (string, error) {
	client := http.Client{}

	body := pk + " " + gammaG2 + " " + partDec

	postReader := strings.NewReader(body)

	resp, err := client.Post("http://127.0.0.1:8080/verifypart", "application/x-www-form-urlencoded", postReader)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", err
	}

	var bodyBytes []byte
	_, err = resp.Body.Read(bodyBytes)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}
