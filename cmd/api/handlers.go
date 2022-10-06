package main

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func (app *application) Handle(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{
		"resp": "ok",
	}

	out, err := json.Marshal(resp)
	if err != nil {
		app.logger.Error("failed to marshal json: ", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
