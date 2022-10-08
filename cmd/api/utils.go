package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const (
	// add more environments
	// handle differences between envs
	develop = "develop"
)

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

// sends requests based on url, method and payload, returns response and error
func sendGetRequest(url string, timeout int) (string, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", 0, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", 0, err
	}

	return string(body), res.StatusCode, nil
}

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
