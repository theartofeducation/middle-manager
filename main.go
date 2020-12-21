package main

import (
	"context"
	"fmt"
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

var (
	port      string
	errorChan chan error
	log       *logrus.Logger
)

func main() {
	log = logrus.New()

	if err := godotenv.Load(); err != nil {
		log.Infoln("could not load env:", err)
	}

	port = os.Getenv("PORT")

	router := mux.NewRouter()

	loadRoutes(router)

	server := &http.Server{
		Addr:         "0.0.0.0:" + port,
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
	log.Infof("Starting server on port %s", port)

	errorChan <- server.ListenAndServe()
}

func handleInterrupt() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT)

	errorChan <- fmt.Errorf("%s", <-ch)
}
