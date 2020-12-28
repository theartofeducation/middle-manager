package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

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
		signature := request.Header.Get("X-Signature")
		body := getBody(request)

		if err := cuClient.VerifySignature(signature, body); err != nil {
			log.Errorln(errors.Wrap(err, "taskStatusUpdatedHandler > verifying signature"))
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		webhook, err := cuClient.GetWebhook(request.Body)
		if err != nil {
			log.Errorln(err)
			writer.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		task, err := cuClient.GetTask(webhook.TaskID)
		if err != nil {
			log.Errorln(err)
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

			client := &http.Client{}
			url := chClient.URL + "/epics"
			req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
			req.Header.Add("Clubhouse-Token", chClient.Token)
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

		writer.WriteHeader(http.StatusNoContent)
	}
}
