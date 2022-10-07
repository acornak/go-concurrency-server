package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// define routes
func (app *application) routes() (r *mux.Router) {
	r = mux.NewRouter()

	// routes
	r.HandleFunc("/v1/api/smart", app.SmartHandler).
		Methods(http.MethodGet).
		// alternative to handle numeric values only in the routes level:
		// Queries("timeout", "{timeout:[0-9]+}")
		Queries("timeout", "{timeout}")

	return
}
