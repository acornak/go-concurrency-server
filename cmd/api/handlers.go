package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (app *application) SmartHandler(w http.ResponseWriter, r *http.Request) {
	var out map[string]string

	start := time.Now()

	timeoutParam := mux.Vars(r)["timeout"]

	// check type of timeout param
	timeout, err := strconv.Atoi(timeoutParam)
	if err != nil {
		app.logger.Error("invalid timeout param: ", zap.Error(err))
		app.errorJson(w, errors.New("invalid timeout parameter: accepts only numbers"))
		return
	}

	// set timeout
	app.timeout = timeout

	// TODO: add logic
	resp, err := app.SendConcurrentRequests()
	if err != nil {
		app.logger.Error("failed to send HTTP request: ", zap.Error(err))
		app.errorJson(w, errors.New("failed to send HTTP request"))
		return
	}

	// check performance
	end := time.Since(start).Milliseconds()

	// setup response
	out = map[string]string{
		"timeout":             fmt.Sprintf("%d ms", timeout),
		"request_performance": fmt.Sprintf("%d ms", end),
		"server_response":     resp,
	}

	// write response
	if err = app.writeJson(w, http.StatusOK, out, ""); err != nil {
		app.logger.Error("failed to marshal json: ", zap.Error(err))
		app.errorJson(w, errors.New("failed to marshal json"))
		return
	}

	// log endpoint performance internally
	// TODO: create logs database for further performance analysis

	app.logger.Info("request performance: ", end, " ms.")
}

func (app *application) SendConcurrentRequests() (resp string, err error) {
	resp, statusCode, err := sendGetRequest(os.Getenv("EXPONEA_URL"))
	if err != nil {
		app.logger.Error("failed to send HTTP request: ", zap.Error(err))
		return
	}

	if statusCode != http.StatusOK {
		// TODO: error chan
		return "", errors.New("server responded with non-OK status")
	} else {
		// TODO: success chan
		app.logger.Info("server response: ", resp)
	}

	return
}
