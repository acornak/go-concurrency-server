package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (app *application) SmartHandler(w http.ResponseWriter, r *http.Request) {
	var resp map[string]string

	start := time.Now()

	timeoutParam := mux.Vars(r)["timeout"]

	// check type of timeout param
	// this is already handled in routes as this endpoint accepts only numeric values for timeout param
	timeout, err := strconv.Atoi(timeoutParam)
	if err != nil {
		app.logger.Error("invalid timeout param: ", zap.Error(err))
		app.errorJson(w, errors.New("invalid timeout parameter: accepts only numbers"))
		return
	}

	// set timeout
	app.timeout = timeout

	// TODO: add logic

	// setup response
	resp = map[string]string{
		"timeout": fmt.Sprintf("%d", timeout),
	}

	// write response
	if err = app.writeJson(w, http.StatusOK, resp, ""); err != nil {
		app.logger.Error("failed to marshal json: ", zap.Error(err))
		app.errorJson(w, errors.New("failed to marshal json"))
		return
	}

	// log endpoint performance internally
	// TODO: create logs database for further performance analysis
	end := time.Since(start).Milliseconds()
	app.logger.Info("request performance: ", end, " ms.")
}
