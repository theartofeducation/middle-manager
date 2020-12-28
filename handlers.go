package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/theartofeducation/middle-manager/clickup"
	"github.com/theartofeducation/middle-manager/clubhouse"
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
			log.Errorln(errors.Wrap(ErrSignatureMismatch, "taskStatusUpdatedHandler > verifying signature"))
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		var webhook clickup.Webhook
		if err := json.NewDecoder(request.Body).Decode(&webhook); err != nil {
			log.Errorln(errors.Wrap(err, "taskStatusUpdatedHandler > decoding webhook"))
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
			log.Errorln(errors.Wrap(err, "taskStatusUpdateHandler > getting task from ClickUp"))
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		if resp.StatusCode != http.StatusOK {
			log.Errorln("taskStatusUpdateHandler > could not get Task from ClickUp with status", resp.StatusCode)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		var task clickup.Task
		if err := json.NewDecoder(resp.Body).Decode((&task)); err != nil {
			log.Errorln(errors.Wrap(err, "taskStatusUpdatedHandler > decoding task"))
			writer.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		if task.Status.Status == clickup.StatusReadyForDevelopment {
			epic := clubhouse.Epic{
				Name:        task.Name,
				Description: task.URL,
			}

			body, err := json.Marshal(epic)
			if err != nil {
				log.Errorln(errors.Wrap(err, "taskStatusUpdateHandler > marshalling create Epic body"))
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}

			client = &http.Client{}
			url := os.Getenv("CLUBHOUSE_API_URL") + "/epics"
			req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
			req.Header.Add("Clubhouse-Token", os.Getenv("CLUBHOUSE_API_TOKEN"))
			req.Header.Add("Content-Type", "application/json")
			clubhouseResponse, err := client.Do(req)

			if err != nil {
				log.Errorln(errors.Wrap(err, "taskStatusUpdateHandler > send request to clubhouse"))
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}

			if clubhouseResponse.StatusCode != http.StatusCreated {
				log.Errorln("taskStatusUpdatedHandler > Clubhouse Epic not created with status", clubhouseResponse.StatusCode)
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}

			log.Infoln("Created Epic:", epic.Name)
		}
	}
}
