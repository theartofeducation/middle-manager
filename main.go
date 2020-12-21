package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
)

var (
	errorChan chan error
	log       *logrus.Logger
)

func main() {
	log = logrus.New()

	if err := godotenv.Load(); err != nil {
		log.Infoln("could not load env:", err)
	}

	router := mux.NewRouter()

	loadRoutes(router)

	server := &http.Server{
		Addr:         "0.0.0.0:" + os.Getenv(("PORT")),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	errorChan = make(chan error, 2)

	go startServer(server)
	go handleInterrupt()

	err := <-errorChan
	err = errors.Wrap(err, "main")
	log.Errorln(err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	server.Shutdown(ctx)
	log.Infoln("shutting down...")
	os.Exit(0)
}

func startServer(server *http.Server) {
	log.Infof("Starting server on %s", server.Addr)

	errorChan <- http.ListenAndServe(server.Addr, handlers.CompressHandler(server.Handler))
}

func handleInterrupt() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT)

	errorChan <- fmt.Errorf("%s", <-ch)
}
