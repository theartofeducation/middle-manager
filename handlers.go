package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/theartofeducation/clickup-go"
	"github.com/theartofeducation/clubhouse-go"
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

		if signature == "" {
			app.log.Errorln(ErrMissingSignature)
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		if body == nil {
			app.log.Errorln(ErrEmptyBody)
			writer.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

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
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		if task.Status.Status == clickup.StatusReadyForDevelopment {
			name := fmt.Sprintf("%s (%s)", task.Name, task.ID)
			epic, err := app.clubhouse.CreateEpic(name, task.URL)
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

func updateTaskStatusHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// validate signature

		// get webhook from clubhouse
		webhook, err := app.clubhouse.ParseWebhook(request.Body)
		if err != nil {
			app.log.Errorln(err)
			// TODO: handler error
			return
		}

		// check webhook for event where an epic was updated to done
		for _, action := range webhook.Actions {
			if epicIsDone(action) {
				update := clickup.UpdateTaskRequest{Status: clickup.StatusAcceptance}

				r := regexp.MustCompile(`((.*))`)
				name := action.Name
				id := r.FindStringSubmatch(name)[1]

				err := app.clickup.UpdateTask(id, update)
				if err != nil {
					app.log.Errorln(err)
					// TODO: handle error
					return
				}

				// log
				task := clickup.Task{}
				app.log.Infof("Task %q moved to acceptance", task.Name)
			}
		}

		writer.WriteHeader(http.StatusNoContent)
	}
}

func epicIsDone(action clubhouse.WebhookAction) bool {
	return action.EntityType == clubhouse.EntityTypeEpic && action.Action == clubhouse.ActionUpdate && action.Changes.State.New == clubhouse.EpicStateDone
}
