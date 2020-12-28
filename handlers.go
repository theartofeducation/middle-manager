package main

import (
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/theartofeducation/middle-manager/clickup"
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
			epic, err := chClient.CreateEpic(task.Name, task.URL)
			if err != nil {
				log.Errorln(err)
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}

			log.Infoln("Created Epic:", epic.Name)
		}

		writer.WriteHeader(http.StatusNoContent)
	}
}
