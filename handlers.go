package main

import (
	"io"
	"net/http"
)

const rootResponse = `{"message": "Locke, I told you I need those TPS reports done by noon today."}`

func rootHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)

		io.WriteString(writer, rootResponse)
	}
}

func taskStatusUpdatedHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		// Verify signature
		writer.WriteHeader(http.StatusOK)

		// TODO: Get Task information
		// TODO: Create Clubhouse Epic
		// TODO: Send Epic to Clubhouse https://clubhouse.io/api/rest/v3/#Create-Epic

		io.WriteString(writer, `{"message": "It works!"}`)
	}
}
