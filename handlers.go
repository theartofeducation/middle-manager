package main

import (
	"io"
	"net/http"
)

func rootHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Type", "application/json")

		io.WriteString(writer, `{"message": "Locke, I told you I need those TPS reports done by noon today."}`)
	}
}
