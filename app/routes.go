package app

import "net/http"

// routes holds all registered routes and universal middleware for the App.
func (app *App) routes() {
	app.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Content-Type", "application/json")

			next.ServeHTTP(writer, request)
		})
	})

	app.router.HandleFunc("/", app.handleIndex()).Methods(http.MethodGet)
	app.router.HandleFunc("/test-error", app.handleTestError()).Methods(http.MethodGet)
}
