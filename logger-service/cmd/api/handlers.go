package main

import (
	"logger/data"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (cfg *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	requestPayload := JSONPayload{}
	cfg.readJSON(w, r, &requestPayload)

	logEntry := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := cfg.Models.LogEntry.Insert(logEntry)
	if err != nil {
		cfg.errorJSON(w, err)
		return
	}

	cfg.writeJSON(w, http.StatusAccepted, jsonResponse{
		Error:   false,
		Message: "logged successfully",
	})

}
