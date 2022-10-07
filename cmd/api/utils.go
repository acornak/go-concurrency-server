package main

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

const (
	// add more environments
	// handle differences between envs
	develop = "develop"
)

// checks if slice a contains string b
func sliceContains(a []string, b string) bool {
	for _, v := range a {
		if v == b {
			return true
		}
	}
	return false
}

// compare if each element of slice a is present in slice b
func compareSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	// this works just fine for smaller slices
	// otherwise this would require different solution due to the overhead
	for _, j := range a {
		if !sliceContains(b, j) {
			return false
		}
	}

	return true
}

// writes successful JSON message
func (app *application) writeJson(w http.ResponseWriter, status int, data interface{}, wrap string) error {
	if wrap != "" {
		wrapper := make(map[string]interface{})
		wrapper[wrap] = data
		data = wrapper
	}

	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err = w.Write(js); err != nil {
		return err
	}

	return nil
}

// writes error JSON message
func (app *application) errorJson(w http.ResponseWriter, err error, status ...int) {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	type jsonError struct {
		Message string `json:"message"`
	}

	errMessage := jsonError{
		Message: err.Error(),
	}

	if err := app.writeJson(w, statusCode, errMessage, "error"); err != nil {
		app.logger.Error("failed to write response: ", zap.Error(err))
	}
}
