package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

func (app *App) writeJSONResponse(writer http.ResponseWriter, statusCode int, v interface{}) {
	writer.WriteHeader(statusCode)

	if err := json.NewEncoder(writer).Encode(v); err != nil {
		app.log.Errorln(errors.Wrap(err, "app.writeJSONResponse"))
	}
}

func (app *App) writeErrorJSONResponse(writer http.ResponseWriter, statusCode int, err error) {
	writer.WriteHeader(statusCode)

	body := fmt.Sprintf(`{"error": %q}`, err.Error())
	if _, err := io.WriteString(writer, body); err != nil {
		app.log.Errorln(errors.Wrap(err, "app.writeErrorJSONResponse"))
	}
}
