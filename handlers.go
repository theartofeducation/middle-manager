package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

const rootResponse = `{"message": "Locke, I told you I need those TPS reports done by noon today."}`

// Custom errors.
var (
	ErrSignatureMismatch = errors.New("Signature mismatch")
)

func rootHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)

		io.WriteString(writer, rootResponse)
	}
}

func taskStatusUpdatedHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusNoContent)

		clickUpSignature := request.Header.Get("X-Signature")
		secret := []byte(os.Getenv("TASK_STATUS_UPDATED_SECRET"))
		body := getBody(request)

		hash := hmac.New(sha256.New, secret)
		hash.Write(body)
		generatedSignature := hex.EncodeToString(hash.Sum(nil))

		if clickUpSignature != generatedSignature {
			log.Errorln(errors.Wrap(ErrSignatureMismatch, "taskStatusUpdatedHandler"))

			writer.WriteHeader(http.StatusUnauthorized)

			return
		}

		var webhook Webhook
		if err := json.NewDecoder(request.Body).Decode(&webhook); err != nil {
			log.Errorln(errors.Wrap(err, "taskStatusUpdatedHandler"))
		}

		log.Infoln("Task ID:", webhook.TaskID)

		// TODO: Get Task information
		// TODO: Create Clubhouse Epic
		// TODO: Send Epic to Clubhouse https://clubhouse.io/api/rest/v3/#Create-Epic
	}
}

// Webhook holds the information for a webhook from ClickUp.
type Webhook struct {
	ID     string `json:"webhook_id"`
	Event  string
	TaskID string `json:"task_id"`
}

func getBody(request *http.Request) (bodyBytes []byte) {
	if request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(request.Body)
		request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	return bodyBytes
}
