package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// define routes
func (app *application) routes() http.Handler {
	r := mux.NewRouter()

	// routes
	r.HandleFunc("/v1/api/smart", app.SendRequestsHandler).
		Methods(http.MethodGet).
		Queries("timeout", "{timeout:[0-9]+}")

	if err := http.ListenAndServe(fmt.Sprintf(":%d", app.config.port), r); err != nil {
		app.logger.Error("failed to serve http: ", zap.Error(err))
	}

	return r
}
