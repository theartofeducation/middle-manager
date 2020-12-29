package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus/hooks/test"
	"github.com/theartofeducation/middle-manager/clickup"
	"github.com/theartofeducation/middle-manager/clubhouse"
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

func Test_taskStatusUpdatedHandler(t *testing.T) {
	t.Run("it creates an Epic when Task is updated to Ready For Development", func(t *testing.T) {
		logger, _ := test.NewNullLogger()
		app = App{
			log:       logger,
			clickup:   clickup.MockClient{},
			clubhouse: clubhouse.MockClient{},
		}

		body := []byte(`{}`)
		request, _ := http.NewRequest(http.MethodPost, "/task-status-updated", bytes.NewBuffer(body))
		request.Header.Add("X-Signature", "abc123")

		response := httptest.NewRecorder()

		handler := http.HandlerFunc(taskStatusUpdatedHandler())
		handler.ServeHTTP(response, request)

		if response.Code != http.StatusNoContent {
			t.Errorf("respones return wrong status: got %d want %d", response.Code, http.StatusNoContent)
		}
	})
}
