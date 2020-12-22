package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

const rootResponse = `{"message": "Locke, I told you I need those TPS reports done by noon today."}`

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
