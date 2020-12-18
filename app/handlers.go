package app

import (
	"errors"
	"net/http"
)

func (app *App) handleIndex() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		app.writeJSONResponse(writer, http.StatusOK, "Hello world!")
	}
}

func (app *App) handleTestError() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		app.writeErrorJSONResponse(writer, http.StatusInternalServerError, errors.New("test error"))
	}
}
