package main

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func loadRoutes(router *mux.Router) {
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Content-Type", "application/json")

			next.ServeHTTP(writer, request)
		})
	})

	router.Handle("/", handlers.CombinedLoggingHandler(os.Stdout, rootHandler())).Methods(http.MethodGet)

	// TODO: Create handler for taskStatusUpdated https://clickup20.docs.apiary.io/#reference/0/webhooks/create-webhook
	// TODO: Get Task information
	// TODO: Create Clubhouse Epic
	// TODO: Send Epic to Clubhouse https://clubhouse.io/api/rest/v3/#Create-Epic
}
