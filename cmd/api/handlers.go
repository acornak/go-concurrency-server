package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (app *application) SendRequestsHandler(w http.ResponseWriter, r *http.Request) {
	var resp map[string]string

	start := time.Now()

	timeoutParam := mux.Vars(r)["timeout"]

	// check type of timeout param
	timeout, err := strconv.Atoi(timeoutParam)
	if err != nil {
		app.logger.Error("invalid timeout param: ", zap.Error(err))
		// TODO: return error
	}

	// set timeout
	app.timeout = timeout

	// TODO: add logic

	resp = map[string]string{
		"timeout": fmt.Sprintf("%d", timeout),
	}

	out, err := json.Marshal(resp)
	if err != nil {
		app.logger.Error("failed to marshal json: ", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(out); err != nil {
		app.logger.Error("failed to write response: ", zap.Error(err))
	}

	end := time.Since(start).Milliseconds()
	app.logger.Info("request performance: ", end, " ms.")
}
