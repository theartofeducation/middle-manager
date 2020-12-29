package main

import (
	"io"
	"net/http"

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

		if err := app.clickup.VerifySignature(signature, body); err != nil {
			app.log.Errorln(err)
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		webhook, err := app.clickup.ParseWebhook(request.Body)
		if err != nil {
			app.log.Errorln(err)
			writer.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		task, err := app.clickup.GetTask(webhook.TaskID)
		if err != nil {
			app.log.Errorln(err)
			writer.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		if task.Status.Status == clickup.StatusReadyForDevelopment {
			epic, err := app.clubhouse.CreateEpic(task.Name, task.URL)
			if err != nil {
				app.log.Errorln(err)
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}

			app.log.Infoln("Created Epic:", epic.Name)
		}

		writer.WriteHeader(http.StatusNoContent)
	}
}
