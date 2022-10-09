package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func init() {
	testApp = application{
		config: config{
			port: 4001,
			env:  develop,
		},
		version: version,
	}
}

func Test_HandleGetRequestError(t *testing.T) {
	testSuccessChan := make(chan string)
	testFailChan := make(chan bool)

	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	logger = zap.New(observedZapCore).Sugar()

	testApp.logger = logger

	testApp.handleGetRequest(testSuccessChan, testFailChan, 1000, func(url string, timeout int) (string, int, error) {
		return "", 0, errors.New("testing error")
	})

	require.Equal(t, 1, observedLogs.Len())
	assert.Equal(t, "failed to send HTTP request: {error 26 0  testing error}", observedLogs.All()[0].Message)
	assert.Equal(t, 0, len(testSuccessChan))
	assert.Equal(t, 0, len(testFailChan))

	close(testSuccessChan)
	close(testFailChan)
}

func Test_HandleGetRequestFail(t *testing.T) {
	testSuccessChan := make(chan string)
	testFailChan := make(chan bool)
	testMessage := "request failed"

	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	logger = zap.New(observedZapCore).Sugar()

	testApp.logger = logger

	go testApp.handleGetRequest(testSuccessChan, testFailChan, 1000, func(url string, timeout int) (string, int, error) {
		return testMessage, http.StatusGatewayTimeout, nil
	})

	assert.Equal(t, 0, len(testFailChan))

	resp := <-testFailChan
	require.Equal(t, 1, observedLogs.Len())
	assert.Equal(t, fmt.Sprintf("server status code: %d response: %s", http.StatusGatewayTimeout, testMessage), observedLogs.All()[0].Message)
	assert.Equal(t, 0, len(testSuccessChan))
	assert.Equal(t, true, resp)

	close(testSuccessChan)
	close(testFailChan)
}

func Test_HandleGetRequestSuccess(t *testing.T) {
	testSuccessChan := make(chan string)
	testFailChan := make(chan bool)
	testMessage := "request successful"

	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	logger = zap.New(observedZapCore).Sugar()

	testApp.logger = logger

	go testApp.handleGetRequest(testSuccessChan, testFailChan, 1000, func(url string, timeout int) (string, int, error) {
		return testMessage, http.StatusOK, nil
	})

	assert.Equal(t, 0, len(testSuccessChan))

	resp := <-testSuccessChan
	require.Equal(t, 1, observedLogs.Len())
	assert.Equal(t, fmt.Sprintf("server status code: %d response: %s", http.StatusOK, testMessage), observedLogs.All()[0].Message)
	assert.Equal(t, 0, len(testFailChan))
	assert.Equal(t, "request successful", resp)

	close(testSuccessChan)
	close(testFailChan)
}

func Test_SmartHandlerInvalidTimeout(t *testing.T) {
	r, err := http.NewRequest("GET", "/v1/api/smart", nil)
	if err != nil {
		t.Error("unable to create new request: ", err)
	}
	q := url.Values{}
	q.Add("timeout", "abcd")
	r.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()

	testApp.SmartHandler(w, r)

	expectedResponse := "{\"error\":{\"message\":\"invalid timeout parameter: accepts only numbers\"}}"

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, []byte(expectedResponse), w.Body.Bytes())
}
