package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_rootHandler(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	handler := http.HandlerFunc(rootHandler())
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("response returned wrong status: got %d want %d", response.Code, http.StatusOK)
	}

	if response.Body.String() != rootResponse {
		t.Errorf("response returned wrong body: got %q want %q", response.Body.String(), rootResponse)
	}
}
