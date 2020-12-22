package main

import (
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

		if !signatureVerified(request) {
			log.Errorln(errors.Wrap(ErrSignatureMismatch, "taskStatusUpdatedHandler > Verify Signature"))

			writer.WriteHeader(http.StatusUnauthorized)

			return
		}

		var webhook Webhook
		if err := json.NewDecoder(request.Body).Decode(&webhook); err != nil {
			log.Errorln(errors.Wrap(err, "taskStatusUpdatedHandler > Decoding Webhook"))

			writer.WriteHeader(http.StatusUnprocessableEntity)

			return
		}

		client := &http.Client{}
		url := os.Getenv("CLICKUP_API_URL") + "/task/" + webhook.TaskID
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Add("Authorization", os.Getenv("CLICKUP_API_KEY"))
		req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)

		if err != nil {
			log.Errorln(errors.Wrap(err, "taskStatusUpdateHandler > ClickUp Get Task"))

			writer.WriteHeader(http.StatusInternalServerError)

			return
		}

		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Errorln(errors.Wrap(errors.New("API error"), "taskStatusUpdateHandler > ClickUp API Response"))

			writer.WriteHeader(http.StatusInternalServerError)

			return
		}

		var task Task
		if err := json.NewDecoder(resp.Body).Decode((&task)); err != nil {
			log.Errorln(errors.Wrap(err, "taskStatusUpdatedHandler > Decoding Task"))

			writer.WriteHeader(http.StatusUnprocessableEntity)

			return
		}

		if task.Status.Status == clickUpStatusReadyForDevelopment {
			// Create Clubhouse Epic
			// Send Epic to Clubhouse https://clubhouse.io/api/rest/v3/#Create-Epic
		}
	}
}
