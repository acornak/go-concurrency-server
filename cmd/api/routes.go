package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// define routes
func (app *application) routes() http.Handler {
	r := mux.NewRouter()

	// routes
	r.HandleFunc("/v1/api/smart", app.Handle).Methods(http.MethodGet)

	http.ListenAndServe(fmt.Sprintf(":%d", app.config.port), r)
	return r
}
