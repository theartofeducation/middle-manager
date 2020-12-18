package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/pkg/errors"

	"github.com/getsentry/sentry-go"

	"github.com/theartofeducation/go-template-repo/app"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/sirupsen/logrus"
)

// env variables
var (
	port string // The HTTP port the server will run on.
	dsn  string // The Sentry DSN.
)

func main() {
	log := logrus.New()

	loadEnvVariables(log)

	args := app.Args{
		Router: mux.NewRouter(),
		Log:    log,
		DSN:    dsn,
	}
	a := app.NewApp(args)

	errorChan := make(chan error, 2)
	go a.StartServer(errorChan, port)
	go handleInterrupt(errorChan)

	err := <-errorChan
	err = errors.Wrap(err, "main")
	sentry.CaptureMessage(err.Error())
	log.Errorln(err)

	a.Close()
}

func loadEnvVariables(log *logrus.Logger) {
	if err := godotenv.Load(); err != nil {
		log.Infoln("could not load env file")
	}

	port = os.Getenv("PORT")
	if strings.TrimSpace(port) == "" {
		log.Fatal("port was not specified")
	}

	dsn = os.Getenv("SENTRY_DSN")
	if strings.TrimSpace(dsn) == "" {
		log.Infoln("Sentry DSN not specified")
	}
}

func handleInterrupt(errorChan chan error) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)
	errorChan <- fmt.Errorf("%s", <-signalChan)
}
