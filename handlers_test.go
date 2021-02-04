package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus/hooks/test"
	"github.com/theartofeducation/clickup-go"
	"github.com/theartofeducation/clubhouse-go"
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

	t.Run("it returns 500 if epic cannot be created", func(t *testing.T) {
		status := clickup.TaskStatus{Status: clickup.StatusReadyForDevelopment}
		task := clickup.Task{Status: status}

		logger, hook := test.NewNullLogger()
		app = App{
			log:       logger,
			clickup:   clickup.MockClient{Task: task},
			clubhouse: clubhouse.MockClient{CreateEpicError: true},
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

		want := fmt.Sprint(clubhouse.ErrTest)
		if hook.LastEntry().Message != want {
			t.Errorf("creation log is incorrect: got %q want %q", hook.LastEntry().Message, want)
		}
	})
}

func Test_updateTaskStatusHandler(t *testing.T) {
	t.Run("it moves a done epic to acceptance", func(t *testing.T) {
		action := clubhouse.WebhookAction{
			Name:       "Test Epic (def456)",
			EntityType: clubhouse.EntityTypeEpic,
			Action:     clubhouse.ActionUpdate,
			Changes: clubhouse.WebhookActionChanges{
				State: clubhouse.WebhookActionState{
					New: clubhouse.EpicStateDone,
				},
			},
		}

		var actions []clubhouse.WebhookAction
		actions = append(actions, action)

		logger, hook := test.NewNullLogger()
		app = App{
			log:     logger,
			clickup: clickup.MockClient{},
			clubhouse: clubhouse.MockClient{
				Webhook: clubhouse.Webhook{
					Actions: actions,
				},
			},
		}

		body := []byte(`{"actions": [{"entity_type": "epic", "action": "update", "changes": {"state": {"new": "done"}}}]}`)
		request, _ := http.NewRequest(http.MethodPost, "/update-task-status", bytes.NewBuffer(body))
		request.Header.Add("Clubhouse-Signature", "abc123")

		response := httptest.NewRecorder()

		handler := http.HandlerFunc(updateTaskStatusHandler())
		handler.ServeHTTP(response, request)

		if response.Code != http.StatusNoContent {
			t.Fatalf("respones return wrong status: got %d want %d", response.Code, http.StatusNoContent)
		}

		want := fmt.Sprintf("Task %q moved to %s", action.Name, clickup.StatusAcceptance)
		if hook.LastEntry().Message != want {
			t.Errorf("unexpected log message: got %q want %q", hook.LastEntry().Message, want)
		}
	})

	t.Run("it moves an in progress epic to in development", func(t *testing.T) {
		action := clubhouse.WebhookAction{
			Name:       "Test Epic (def456)",
			EntityType: clubhouse.EntityTypeEpic,
			Action:     clubhouse.ActionUpdate,
			Changes: clubhouse.WebhookActionChanges{
				State: clubhouse.WebhookActionState{
					New: clubhouse.EpicStateInProgress,
				},
			},
		}

		var actions []clubhouse.WebhookAction
		actions = append(actions, action)

		logger, hook := test.NewNullLogger()
		app = App{
			log:     logger,
			clickup: clickup.MockClient{},
			clubhouse: clubhouse.MockClient{
				Webhook: clubhouse.Webhook{
					Actions: actions,
				},
			},
		}

		body := []byte(`{"actions": [{"entity_type": "epic", "action": "update", "changes": {"state": {"new": "in progress"}}}]}`)
		request, _ := http.NewRequest(http.MethodPost, "/update-task-status", bytes.NewBuffer(body))
		request.Header.Add("Clubhouse-Signature", "abc123")

		response := httptest.NewRecorder()

		handler := http.HandlerFunc(updateTaskStatusHandler())
		handler.ServeHTTP(response, request)

		if response.Code != http.StatusNoContent {
			t.Fatalf("respones return wrong status: got %d want %d", response.Code, http.StatusNoContent)
		}

		want := fmt.Sprintf("Task %q moved to %s", action.Name, clickup.StatusInDevelopmentClubhouse)
		if hook.LastEntry().Message != want {
			t.Errorf("unexpected log message: got %q want %q", hook.LastEntry().Message, want)
		}
	})

	t.Run("it handles empty signature", func(t *testing.T) {
		logger, hook := test.NewNullLogger()
		app = App{
			log:       logger,
			clickup:   clickup.MockClient{},
			clubhouse: clubhouse.MockClient{},
		}

		body := []byte(`{"actions": [{"entity_type": "epic", "action": "update", "changes": {"state": {"new": "done"}}}]}`)
		request, _ := http.NewRequest(http.MethodPost, "/update-task-status", bytes.NewBuffer(body))
		request.Header.Add("Clubhouse-Signature", "")

		response := httptest.NewRecorder()

		handler := http.HandlerFunc(updateTaskStatusHandler())
		handler.ServeHTTP(response, request)

		if response.Code != http.StatusUnauthorized {
			t.Fatalf("respones return wrong status: got %d want %d", response.Code, http.StatusUnauthorized)
		}

		if hook.LastEntry().Message != ErrMissingSignature.Error() {
			t.Errorf("received wrong error message: got %q want %q", hook.LastEntry().Message, ErrMissingSignature)
		}
	})

	t.Run("it handles empty body", func(t *testing.T) {
		logger, hook := test.NewNullLogger()
		app = App{
			log:       logger,
			clickup:   clickup.MockClient{},
			clubhouse: clubhouse.MockClient{},
		}

		request, _ := http.NewRequest(http.MethodPost, "/update-task-status", nil)
		request.Header.Add("Clubhouse-Signature", "abc123")

		response := httptest.NewRecorder()

		handler := http.HandlerFunc(updateTaskStatusHandler())
		handler.ServeHTTP(response, request)

		if response.Code != http.StatusUnprocessableEntity {
			t.Fatalf("respones return wrong status: got %d want %d", response.Code, http.StatusUnprocessableEntity)
		}

		if hook.LastEntry().Message != ErrEmptyBody.Error() {
			t.Errorf("received wrong error message: got %q want %q", hook.LastEntry().Message, ErrEmptyBody)
		}
	})

	t.Run("it handles signature mismatch", func(t *testing.T) {
		logger, hook := test.NewNullLogger()
		app = App{
			log:       logger,
			clickup:   clickup.MockClient{},
			clubhouse: clubhouse.MockClient{VerifySignatureError: true},
		}

		body := []byte(`{"actions": [{"entity_type": "epic", "action": "update", "changes": {"state": {"new": "done"}}}]}`)
		request, _ := http.NewRequest(http.MethodPost, "/update-task-status", bytes.NewBuffer(body))
		request.Header.Add("Clubhouse-Signature", "abc123")

		response := httptest.NewRecorder()

		handler := http.HandlerFunc(updateTaskStatusHandler())
		handler.ServeHTTP(response, request)

		if response.Code != http.StatusUnauthorized {
			t.Fatalf("respones return wrong status: got %d want %d", response.Code, http.StatusUnauthorized)
		}

		if hook.LastEntry().Message != clubhouse.ErrSignatureMismatch.Error() {
			t.Errorf("received wrong error message: got %q want %q", hook.LastEntry().Message, clubhouse.ErrSignatureMismatch)
		}
	})

	t.Run("it handles parse webhook error", func(t *testing.T) {
		logger, hook := test.NewNullLogger()
		app = App{
			log:       logger,
			clickup:   clickup.MockClient{},
			clubhouse: clubhouse.MockClient{ParseWebhookError: true},
		}

		body := []byte(`{"actions": [{"entity_type": "epic", "action": "update", "changes": {"state": {"new": "done"}}}]}`)
		request, _ := http.NewRequest(http.MethodPost, "/update-task-status", bytes.NewBuffer(body))
		request.Header.Add("Clubhouse-Signature", "abc123")

		response := httptest.NewRecorder()

		handler := http.HandlerFunc(updateTaskStatusHandler())
		handler.ServeHTTP(response, request)

		if response.Code != http.StatusUnprocessableEntity {
			t.Fatalf("respones return wrong status: got %d want %d", response.Code, http.StatusUnprocessableEntity)
		}

		if hook.LastEntry().Message != "Could not parse Webhook body: Test error" {
			t.Errorf("received wrong error message: got %q want %q", hook.LastEntry().Message, "Could not parse Webhook body: Test error")
		}
	})

	t.Run("it handles epic without Task ID", func(t *testing.T) {
		action := clubhouse.WebhookAction{
			Name:       "Test Epic",
			EntityType: clubhouse.EntityTypeEpic,
			Action:     clubhouse.ActionUpdate,
			Changes: clubhouse.WebhookActionChanges{
				State: clubhouse.WebhookActionState{
					New: clubhouse.EpicStateDone,
				},
			},
		}

		var actions []clubhouse.WebhookAction
		actions = append(actions, action)

		logger, hook := test.NewNullLogger()
		app = App{
			log:     logger,
			clickup: clickup.MockClient{},
			clubhouse: clubhouse.MockClient{
				Webhook: clubhouse.Webhook{
					Actions: actions,
				},
			},
		}

		body := []byte(`{"actions": [{"entity_type": "epic", "action": "update", "changes": {"state": {"new": "done"}}}]}`)
		request, _ := http.NewRequest(http.MethodPost, "/update-task-status", bytes.NewBuffer(body))
		request.Header.Add("Clubhouse-Signature", "abc123")

		response := httptest.NewRecorder()

		handler := http.HandlerFunc(updateTaskStatusHandler())
		handler.ServeHTTP(response, request)

		if response.Code != http.StatusNoContent {
			t.Fatalf("respones return wrong status: got %d want %d", response.Code, http.StatusNoContent)
		}

		if hook.LastEntry() != nil {
			t.Errorf("received log message when not expecting one: %q", hook.LastEntry().Message)
		}
	})

	t.Run("it handles task update error", func(t *testing.T) {
		action := clubhouse.WebhookAction{
			Name:       "Test Epic (def456)",
			EntityType: clubhouse.EntityTypeEpic,
			Action:     clubhouse.ActionUpdate,
			Changes: clubhouse.WebhookActionChanges{
				State: clubhouse.WebhookActionState{
					New: clubhouse.EpicStateDone,
				},
			},
		}

		var actions []clubhouse.WebhookAction
		actions = append(actions, action)

		logger, hook := test.NewNullLogger()
		app = App{
			log:     logger,
			clickup: clickup.MockClient{UpdateTaskError: true},
			clubhouse: clubhouse.MockClient{
				Webhook: clubhouse.Webhook{
					Actions: actions,
				},
			},
		}

		body := []byte(`{"actions": [{"entity_type": "epic", "action": "update", "changes": {"state": {"new": "done"}}}]}`)
		request, _ := http.NewRequest(http.MethodPost, "/update-task-status", bytes.NewBuffer(body))
		request.Header.Add("Clubhouse-Signature", "abc123")

		response := httptest.NewRecorder()

		handler := http.HandlerFunc(updateTaskStatusHandler())
		handler.ServeHTTP(response, request)

		if response.Code != http.StatusInternalServerError {
			t.Fatalf("respones return wrong status: got %d want %d", response.Code, http.StatusInternalServerError)
		}

		if hook.LastEntry().Message != clickup.ErrTest.Error() {
			t.Errorf("received wrong error message: got %q want %q", hook.LastEntry().Message, clickup.ErrTest)
		}
	})
}
