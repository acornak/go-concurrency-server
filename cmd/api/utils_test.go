package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
)

func Test_SameSlices(t *testing.T) {
	a := []string{"foo", "bar", "sufurki"}
	b := []string{"foo", "bar", "sufurki"}

	if !compareSlices(a, b) {
		t.Error("slices should be equal")
	}
}

func Test_DifferentLengthSlices(t *testing.T) {
	a := []string{"foo", "bar", "sufurki"}
	b := []string{"foo", "bar"}

	if compareSlices(a, b) {
		t.Error("slices should be different")
	}
}

func Test_DifferentSlicesSameLength(t *testing.T) {
	a := []string{"foo", "bar", "sufurki"}
	b := []string{"foo", "bar", "not-sufurki"}

	if compareSlices(a, b) {
		t.Error("slices should be different")
	}
}

func Test_WrapperJsonStatus(t *testing.T) {
	statuses := []int{
		http.StatusOK,
		http.StatusBadGateway,
		http.StatusBadGateway,
		http.StatusForbidden,
	}

	for _, stat := range statuses {
		rr := httptest.NewRecorder()
		err := testApp.writeJson(rr, stat, "", "")
		if err != nil {
			t.Error("test failed: ", zap.Error(err))
		}

		if rr.Code != stat {
			t.Errorf("expected status code of %d, but got %d", stat, rr.Code)
		}
	}
}

func Test_WrapperJsonNoWrapper(t *testing.T) {
	rr := httptest.NewRecorder()
	payload := map[string]string{
		"test": "successful",
	}

	err := testApp.writeJson(rr, http.StatusOK, payload, "")
	if err != nil {
		t.Error("test failed: ", zap.Error(err))
	}

	body, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Error("failed to read body: ", err)
	}

	expected, err := json.Marshal(payload)
	if err != nil {
		t.Error("failed to marshall json: ", err)
	}

	if string(body) != string(expected) {
		t.Errorf("expected response %s, but got %s", expected, body)
	}
}

func Test_WrapperJsonWrapper(t *testing.T) {
	rr := httptest.NewRecorder()
	payload := map[string]string{
		"test": "successful",
	}

	err := testApp.writeJson(rr, http.StatusOK, payload, "wrapper")
	if err != nil {
		t.Error("test failed: ", zap.Error(err))
	}

	body, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Error("failed to read body: ", err)
	}

	expected := "{\"wrapper\":{\"test\":\"successful\"}}"

	if string(body) != expected {
		t.Errorf("expected response %s, but got %s", expected, body)
	}
}

func Test_ErrorJsonStatuses(t *testing.T) {
	statuses := []int{
		http.StatusBadGateway,
		http.StatusBadGateway,
		http.StatusForbidden,
	}

	for _, stat := range statuses {
		rr := httptest.NewRecorder()
		testApp.errorJson(rr, errors.New("testing error"), stat)

		if rr.Code != stat {
			t.Errorf("expected status code of %d, but got %d", stat, rr.Code)
		}
	}
}

func Test_ErrorJsonNoStatus(t *testing.T) {
	rr := httptest.NewRecorder()
	testApp.errorJson(rr, errors.New("testing error"))

	body, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Error("failed to read body: ", err)
	}

	expected := "{\"error\":{\"message\":\"testing error\"}}"

	if string(body) != expected {
		t.Errorf("expected response %s, but got %s", expected, body)
	}

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status code of 400, but got %d", rr.Code)
	}
}
