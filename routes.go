package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func loadRoutes(router *mux.Router) {
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Content-Type", "application/json")

			next.ServeHTTP(writer, request)
		})
	})

	router.HandleFunc("/", rootHandler()).Methods(http.MethodGet)
	router.HandleFunc("/task-status-updated", taskStatusUpdatedHandler()).Methods(http.MethodPost)
	router.HandleFunc("/update-task-status", updateTaskStatusHandler()).Methods(http.MethodPost)
}
