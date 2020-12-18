package apptests

import (
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/theartofeducation/go-template-repo/app"
)

func TestNewApp(t *testing.T) {
	t.Run("it returns a new instance of an App", func(t *testing.T) {
		logger, hook := test.NewNullLogger()

		args := app.Args{
			Router: mux.NewRouter(),
			Log:    logger,
		}

		_ = app.NewApp(args)

		if hook.LastEntry() != nil {
			t.Errorf("did not expect a log message: got %q", hook.LastEntry().Message)
		}
	})
}
