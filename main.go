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

	"github.com/theartofeducation/middle-manager/clickup"
	"github.com/theartofeducation/middle-manager/clubhouse"
)

var (
	errorChan chan error
	app       App
)

// App holds the application dependencies.
type App struct {
	log       *logrus.Logger
	clickup   clickup.CUClient
	clubhouse clubhouse.CHClient
}

func main() {
	app = App{}

	app.log = logrus.New()

	if err := godotenv.Load(); err != nil {
		app.log.Infoln("could not load env:", err)
	}

	app.clickup = clickup.NewClient(os.Getenv("CLICKUP_API_KEY"), os.Getenv("TASK_STATUS_UPDATED_SECRET"))
	app.clubhouse = clubhouse.NewClient(os.Getenv("CLUBHOUSE_API_TOKEN"))

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
	app.log.Errorln(err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	server.Shutdown(ctx)
	app.log.Infoln("shutting down...")
	os.Exit(0)
}

func startServer(server *http.Server) {
	app.log.Infof("Starting server on %s", server.Addr)

	errorChan <- http.ListenAndServe(
		server.Addr,
		handlers.CompressHandler(handlers.CombinedLoggingHandler(os.Stdout, server.Handler)),
	)
}

func handleInterrupt() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT)

	errorChan <- fmt.Errorf("%s", <-ch)
}
