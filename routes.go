package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func loadRoutes(router *mux.Router) {
	router.HandleFunc("/", rootHandler()).Methods(http.MethodGet)
}
