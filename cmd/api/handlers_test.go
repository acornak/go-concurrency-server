package main

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func init() {
	testApp = application{
		config: config{
			port: 1234,
			env:  develop,
		},
		version: version,
	}
}

func Test_SmartHandlerError(t *testing.T) {
	testSuccessChan := make(chan string)
	testFailChan := make(chan failChanStruct)

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

func Test_SmartHandlerFail(t *testing.T) {
	testSuccessChan := make(chan string)
	testFailChan := make(chan failChanStruct)
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
	assert.Equal(t, failChanStruct{
		status:  http.StatusGatewayTimeout,
		message: "request failed",
	}, resp)

	close(testSuccessChan)
	close(testFailChan)
}

func Test_SmartHandlerSuccess(t *testing.T) {
	testSuccessChan := make(chan string)
	testFailChan := make(chan failChanStruct)
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
