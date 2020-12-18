package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

// TODO: add proper logging

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", rootHandler()).Methods(http.MethodGet)

	server := &http.Server{
		Addr:         "0.0.0.0:8080", // TODO: move port to env
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	go serve(server)

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)
	<-interruptChan

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	server.Shutdown(ctx)
	fmt.Println("shutting down...")
	os.Exit(0)

	// Create handler for taskStatusUpdated https://clickup20.docs.apiary.io/#reference/0/webhooks/create-webhook
	// Get Task information
	// Create Clubhouse Epic
	// Send Epic to Clubhouse https://clubhouse.io/api/rest/v3/#Create-Epic
}

func rootHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Type", "application/json")

		io.WriteString(writer, `{"message": "Locke, I told you I need those TPS reports done by noon today."}`)
	}
}

func serve(server *http.Server) {
	fmt.Println("starting server on port 8080...") // TODO: move port to env

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("failed to start the server: ", err) // TODO: error being triggered on shutdown
	}
}
