package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"net/http"
)

func getBody(request *http.Request) (bodyBytes []byte) {
	if request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(request.Body)
		request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	return bodyBytes
}

func signatureVerified(request *http.Request) bool {
	clickUpSignature := request.Header.Get("X-Signature")
	secret := []byte(cuClient.TaskStatusUpdatedSecret)

	body := getBody(request)

	hash := hmac.New(sha256.New, secret)
	hash.Write(body)
	generatedSignature := hex.EncodeToString(hash.Sum(nil))

	return clickUpSignature == generatedSignature
}
