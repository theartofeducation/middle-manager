package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
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
			log.Errorln(errors.Wrap(ErrSignatureMismatch, "taskStatusUpdatedHandler"))

			writer.WriteHeader(http.StatusUnauthorized)

			return
		}

		var webhook Webhook
		if err := json.NewDecoder(request.Body).Decode(&webhook); err != nil {
			log.Errorln(errors.Wrap(err, "taskStatusUpdatedHandler"))

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
		respBody, _ := ioutil.ReadAll(resp.Body)

		log.Infoln("Response status:", resp.Status)
		log.Infoln("Response body:", string(respBody))
		// TODO: unmarshal Task

		// TODO: Create Clubhouse Epic
		// TODO: Send Epic to Clubhouse https://clubhouse.io/api/rest/v3/#Create-Epic
	}
}
