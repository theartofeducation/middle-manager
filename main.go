package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()

	godotenv.Load()

	router := mux.NewRouter()

	router.HandleFunc("/", rootHandler()).Methods(http.MethodGet)

	server := &http.Server{
		Addr:         "0.0.0.0:" + os.Getenv("PORT"),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	go serve(server, log)

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)
	<-interruptChan

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	server.Shutdown(ctx)
	log.Info("shutting down...")
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

func serve(server *http.Server, log *logrus.Logger) {
	log.Infof("starting server on port %s...\n", os.Getenv("PORT"))

	if err := server.ListenAndServe(); err != nil {
		log.Error("failed to start the server: ", err) // TODO: error being triggered on shutdown
	}
}
