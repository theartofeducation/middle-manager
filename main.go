package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()

	godotenv.Load()

	router := mux.NewRouter()

	router.HandleFunc("/", rootHandler()).Methods(http.MethodGet)

	// Create handler for taskStatusUpdated https://clickup20.docs.apiary.io/#reference/0/webhooks/create-webhook
	// Get Task information
	// Create Clubhouse Epic
	// Send Epic to Clubhouse https://clubhouse.io/api/rest/v3/#Create-Epic

	server := &http.Server{
		Addr:         "0.0.0.0:" + os.Getenv("PORT"),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	errorChan := make(chan error, 2)

	go startServer(server, log, errorChan)
	go handleInterrupt(errorChan)

	err := <-errorChan
	err = errors.Wrap(err, "main")
	log.Errorln(err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	server.Shutdown(ctx)
	log.Infoln("shutting down...")
	os.Exit(0)
}

func rootHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Type", "application/json")

		io.WriteString(writer, `{"message": "Locke, I told you I need those TPS reports done by noon today."}`)
	}
}

func startServer(server *http.Server, log *logrus.Logger, errorChan chan error) {
	log.Infof("Starting server on port %s", os.Getenv("PORT"))

	errorChan <- server.ListenAndServe()
}

func handleInterrupt(errorChan chan error) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT)

	errorChan <- fmt.Errorf("%s", <-ch)
}
