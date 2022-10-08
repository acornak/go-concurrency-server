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

type failChanStruct struct {
	status  int
	message string
}

type SmartResponseMessage struct {
	Result      string `json:"result"`
	Performance string `json:"performance"`
	Timeout     string `json:"timeout"`
}

type SmartResponse struct {
	Status  int                  `json:"status"`
	Message SmartResponseMessage `json:"message"`
}

type requestSender func(url string, timeout int) (string, int, error)

// handle `/v1/api/smart` endpoint
func (app *application) SmartHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	out := &SmartResponse{}
	status := http.StatusOK

	successChan := make(chan string)
	failChan := make(chan failChanStruct)
	doneChan := make(chan bool, 1)

	timeoutParam := mux.Vars(r)["timeout"]

	// check type of timeout param
	timeout, err := strconv.Atoi(timeoutParam)
	if err != nil {
		app.logger.Error("invalid timeout param: ", zap.Error(err))
		app.errorJson(w, errors.New("invalid timeout parameter: accepts only numbers"))
		return
	}

	out.Message.Timeout = fmt.Sprintf("%d ms", timeout)

	// try first response
	go app.handleGetRequest(successChan, failChan, timeout, sendGetRequest)

	// handle response
	go func() {
		for {
			select {
			// check if request was successful
			case out.Message.Result = <-successChan:
				doneChan <- true
				close(successChan)
				return
			// check if request was not successful
			case resp := <-failChan:
				doneChan <- true
				status = resp.status
				out.Message.Result = resp.message
				close(failChan)
				return
			// check if server responded within specified timeout
			case <-time.After(time.Duration(timeout) * time.Millisecond):
				status = http.StatusGatewayTimeout
				out.Message.Result = "server did not respond within specified timeout"
				doneChan <- true
				return
			// check if server responded within 300 ms
			case <-time.After(300 * time.Millisecond):
				app.logger.Info("server did not respond within 300 ms")
				status = http.StatusGatewayTimeout
				out.Message.Result = "server did not respond within 300 ms"
				doneChan <- true
				return
			}
		}
	}()

	<-doneChan
	close(doneChan)

	// check performance
	end := time.Since(start).Milliseconds()
	out.Status = status
	out.Message.Performance = fmt.Sprintf("%d ms", end)

	// write response
	if err = app.writeJson(w, status, out, ""); err != nil {
		app.logger.Error("failed to marshal json: ", zap.Error(err))
		app.errorJson(w, errors.New("failed to marshal json"))
		return
	}

	// POTENTIAL PRODUCTION TODOS:
	// log endpoint performance internally
	// create logs database for further performance analysis

	app.logger.Info("request performance: ", end, " ms.")
}

func (app *application) handleGetRequest(successChan chan string, failChan chan failChanStruct, timeout int, getReq requestSender) {
	resp, statusCode, err := getReq(os.Getenv("EXPONEA_URL"), timeout)
	if err != nil {
		app.logger.Error("failed to send HTTP request: ", zap.Error(err))
		return
	}
	app.logger.Infof("server status code: %d response: %s", statusCode, resp)

	if statusCode == http.StatusOK {
		successChan <- resp
	} else {
		failChan <- failChanStruct{
			status:  statusCode,
			message: resp,
		}
	}
}
