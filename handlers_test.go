package main

import (
	"bytes"
	"fmt"
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
		status := clickup.TaskStatus{Status: clickup.StatusReadyForDevelopment}
		task := clickup.Task{Status: status}

		epic := clubhouse.Epic{Name: "Test Epic"}

		logger, hook := test.NewNullLogger()
		app = App{
			log:       logger,
			clickup:   clickup.MockClient{Task: task},
			clubhouse: clubhouse.MockClient{Epic: epic},
		}

		body := []byte(`{"webhook_id": "def456", "event": "taskStatusUpdated", "task_id": "test1"}`)
		request, _ := http.NewRequest(http.MethodPost, "/task-status-updated", bytes.NewBuffer(body))
		request.Header.Add("X-Signature", "abc123")

		response := httptest.NewRecorder()

		handler := http.HandlerFunc(taskStatusUpdatedHandler())
		handler.ServeHTTP(response, request)

		if response.Code != http.StatusNoContent {
			t.Errorf("respones return wrong status: got %d want %d", response.Code, http.StatusNoContent)
		}

		want := fmt.Sprint("Created Epic: ", epic.Name)
		if hook.LastEntry().Message != want {
			t.Errorf("creation log is incorrect: got %q want %q", hook.LastEntry().Message, want)
		}
	})

	t.Run("it returns 401 if signature is not present", func(t *testing.T) {
		logger, hook := test.NewNullLogger()
		app = App{
			log:       logger,
			clickup:   clickup.MockClient{},
			clubhouse: clubhouse.MockClient{},
		}

		body := []byte(`{"webhook_id": "def456", "event": "taskStatusUpdated", "task_id": "test1"}`)
		request, _ := http.NewRequest(http.MethodPost, "/task-status-updated", bytes.NewBuffer(body))

		response := httptest.NewRecorder()

		handler := http.HandlerFunc(taskStatusUpdatedHandler())
		handler.ServeHTTP(response, request)

		if response.Code != http.StatusUnauthorized {
			t.Errorf("respones return wrong status: got %d want %d", response.Code, http.StatusUnauthorized)
		}

		want := fmt.Sprint(ErrMissingSignature)
		if hook.LastEntry().Message != want {
			t.Errorf("creation log is incorrect: got %q want %q", hook.LastEntry().Message, want)
		}
	})

	t.Run("it returns 422 if body is empty", func(t *testing.T) {
		logger, hook := test.NewNullLogger()
		app = App{
			log:       logger,
			clickup:   clickup.MockClient{},
			clubhouse: clubhouse.MockClient{},
		}

		request, _ := http.NewRequest(http.MethodPost, "/task-status-updated", nil)
		request.Header.Add("X-Signature", "abc123")

		response := httptest.NewRecorder()

		handler := http.HandlerFunc(taskStatusUpdatedHandler())
		handler.ServeHTTP(response, request)

		if response.Code != http.StatusUnprocessableEntity {
			t.Errorf("respones return wrong status: got %d want %d", response.Code, http.StatusUnprocessableEntity)
		}

		want := fmt.Sprint(ErrEmptyBody)
		if hook.LastEntry().Message != want {
			t.Errorf("creation log is incorrect: got %q want %q", hook.LastEntry().Message, want)
		}
	})

	t.Run("it returns 401 for unverified signature", func(t *testing.T) {
		logger, hook := test.NewNullLogger()
		app = App{
			log:       logger,
			clickup:   clickup.MockClient{VerifySignatureError: true},
			clubhouse: clubhouse.MockClient{},
		}

		body := []byte(`{"webhook_id": "def456", "event": "taskStatusUpdated", "task_id": "test1"}`)
		request, _ := http.NewRequest(http.MethodPost, "/task-status-updated", bytes.NewBuffer(body))
		request.Header.Add("X-Signature", "abc123")

		response := httptest.NewRecorder()

		handler := http.HandlerFunc(taskStatusUpdatedHandler())
		handler.ServeHTTP(response, request)

		if response.Code != http.StatusUnauthorized {
			t.Errorf("respones return wrong status: got %d want %d", response.Code, http.StatusUnauthorized)
		}

		want := fmt.Sprint(clickup.ErrSignatureMismatch)
		if hook.LastEntry().Message != want {
			t.Errorf("creation log is incorrect: got %q want %q", hook.LastEntry().Message, want)
		}
	})

	t.Run("it returns 422 if webhook cannot be parsed", func(t *testing.T) {
		logger, hook := test.NewNullLogger()
		app = App{
			log:       logger,
			clickup:   clickup.MockClient{ParseWebhookError: true},
			clubhouse: clubhouse.MockClient{},
		}

		body := []byte(`{"webhook_id": "def456", "event": "taskStatusUpdated", "task_id": "test1"}`)
		request, _ := http.NewRequest(http.MethodPost, "/task-status-updated", bytes.NewBuffer(body))
		request.Header.Add("X-Signature", "abc123")

		response := httptest.NewRecorder()

		handler := http.HandlerFunc(taskStatusUpdatedHandler())
		handler.ServeHTTP(response, request)

		if response.Code != http.StatusUnprocessableEntity {
			t.Errorf("respones return wrong status: got %d want %d", response.Code, http.StatusUnprocessableEntity)
		}

		want := fmt.Sprint("Could not parse Webhook body: ", clickup.ErrTest)
		if hook.LastEntry().Message != want {
			t.Errorf("creation log is incorrect: got %q want %q", hook.LastEntry().Message, want)
		}
	})

	t.Run("it returns 500 if task cannot be fetched", func(t *testing.T) {
		logger, hook := test.NewNullLogger()
		app = App{
			log:       logger,
			clickup:   clickup.MockClient{GetTaskError: true},
			clubhouse: clubhouse.MockClient{},
		}

		body := []byte(`{"webhook_id": "def456", "event": "taskStatusUpdated", "task_id": "test1"}`)
		request, _ := http.NewRequest(http.MethodPost, "/task-status-updated", bytes.NewBuffer(body))
		request.Header.Add("X-Signature", "abc123")

		response := httptest.NewRecorder()

		handler := http.HandlerFunc(taskStatusUpdatedHandler())
		handler.ServeHTTP(response, request)

		if response.Code != http.StatusInternalServerError {
			t.Errorf("respones return wrong status: got %d want %d", response.Code, http.StatusInternalServerError)
		}

		want := fmt.Sprint(clickup.ErrTest)
		if hook.LastEntry().Message != want {
			t.Errorf("creation log is incorrect: got %q want %q", hook.LastEntry().Message, want)
		}
	})
}
