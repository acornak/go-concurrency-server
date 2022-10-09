package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type SmartResponseMessage struct {
	Result   string `json:"result"`
	Duration string `json:"duration"`
	Timeout  string `json:"timeout"`
}

type SmartResponse struct {
	Status  int                  `json:"status"`
	Message SmartResponseMessage `json:"message"`
}

type requestSender func(url string, timeout int) (string, int, error)

// handle `/v1/api/smart` endpoint
func (app *application) SmartHandler(w http.ResponseWriter, r *http.Request) {
	reqID := uuid.New()
	app.logger.Infof("---- started handling requestID %s ----", reqID)
	start := time.Now()

	out := &SmartResponse{}
	status := http.StatusOK

	successChan := make(chan string)
	failChan := make(chan bool, 1)
	failChanInternal := make(chan bool, 3)
	doneChan := make(chan bool)
	concurChan := make(chan bool, 1)

	timeoutReq := 300
	timeoutReqTimer := time.After(time.Duration(timeoutReq) * time.Millisecond)

	timeoutParam := mux.Vars(r)["timeout"]

	// check type of timeout param
	timeoutUser, err := strconv.Atoi(timeoutParam)
	if err != nil {
		app.logger.Errorf("invalid timeout param: ", zap.Error(err))
		app.errorJson(w, errors.New("invalid timeout parameter: accepts only numbers"))
		return
	}

	out.Message.Timeout = fmt.Sprintf("%d ms", timeoutUser)

	// try first request
	go app.handleGetRequest(successChan, failChan, timeoutUser, sendGetRequest)

	go func() {
		for {
			select {
			// check the first successful request
			case out.Message.Result = <-successChan:
				app.logger.Info("request response OK! requestID: ", reqID)
				doneChan <- true
			// check if all requests were not successful
			case v := <-failChan:
				failChanInternal <- v
				if len(failChanInternal) == cap(failChanInternal) {
					app.logger.Info("all 3 requests failed. requestID: ", reqID)
					status = http.StatusBadGateway
					out.Message.Result = "all 3 requests failed"
					doneChan <- true
				}
			// check if server responded within 300 ms
			case <-timeoutReqTimer:
				app.logger.Info("server did not respond successfully within 300 ms, firing up another 2 requests. requestID: ", reqID)
				app.logger.Info("current time elapsed: ", time.Since(start).Milliseconds())
				concurChan <- true
			// check if there is need to start 2 concurrent requests
			case <-concurChan:
				app.handleConcurrentRequests(successChan, failChan, timeoutUser-timeoutReq, sendGetRequest)
			}
		}
	}()

	// listen for results
	select {
	// check if requests are done successfully
	case <-doneChan:
	// check if server did not respond within specified timeout
	case <-time.After(time.Duration(timeoutUser) * time.Millisecond):
		app.logger.Info("server did not respond within specified timeout. requestID: ", reqID)
		status = http.StatusGatewayTimeout
		out.Message.Result = "server did not respond within specified timeout"
	}

	// check performance
	end := time.Since(start).Milliseconds()
	out.Status = status
	out.Message.Duration = fmt.Sprintf("%d ms", end)

	// write response
	if err = app.writeJson(w, status, out, ""); err != nil {
		app.logger.Error("failed to marshal json: ", zap.Error(err))
		app.errorJson(w, errors.New("failed to marshal json"))
		return
	}

	// POTENTIAL PRODUCTION TODOS:
	// log endpoint performance internally
	// create logs database for further performance analysis

	app.logger.Infof("requestID: %s, status: %d response: %s duration: %s timeout: %s", reqID, out.Status, out.Message.Result, out.Message.Duration, out.Message.Timeout)
	app.logger.Infof("---- requestID %s completed! ----", reqID)
}

// send GET request and write to channel
func (app *application) handleGetRequest(successChan chan string, failChan chan bool, timeout int, getReq requestSender) {
	resp, statusCode, err := getReq(os.Getenv("EXPONEA_URL"), timeout)
	if err != nil {
		app.logger.Error("failed to send HTTP request: ", zap.Error(err))
		return
	}
	app.logger.Infof("server status code: %d response: %s", statusCode, resp)

	if statusCode == http.StatusOK {
		successChan <- resp
	} else {
		failChan <- true
	}
}

// fire up 2 go routines to send additional requests (further logic to be placed if needed)
func (app *application) handleConcurrentRequests(successChan chan string, failChan chan bool, timeout int, getReq requestSender) {
	for i := 0; i < 2; i++ {
		go app.handleGetRequest(successChan, failChan, timeout, sendGetRequest)
	}
}
