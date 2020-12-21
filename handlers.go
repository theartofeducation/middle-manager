package main

import (
	"encoding/json"
	"io"
	"net/http"

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
		// TODO: filter to list https://app.clickup.com/2324163/v/li/52581966
		writer.WriteHeader(http.StatusNoContent)

		// TODO: Verify signature
		signature := request.Header.Get("X-Signature")
		log.Infoln("Signature: ", signature)

		// key := os.Getenv("TASK_STATUS_UPDATED_SECRET")

		var webhook Webhook
		if err := json.NewDecoder(request.Body).Decode(&webhook); err != nil {
			log.Errorln(errors.Wrap(err, "taskStatusUpdatedHandler"))
		}

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
	// history_items
}
