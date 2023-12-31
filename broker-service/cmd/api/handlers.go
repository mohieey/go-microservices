package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (cfg *Config) Broker(w http.ResponseWriter, r *http.Request) {

	payload := jsonResponse{
		Error:   false,
		Message: "hit the broker",
	}

	cfg.writeJSON(w, http.StatusOK, payload)
}

func (cfg *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := cfg.readJSON(w, r, &requestPayload)
	if err != nil {
		cfg.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		cfg.authenticate(w, requestPayload.Auth)
	case "log":
		cfg.logItem(w, requestPayload.Log)
	default:
		cfg.errorJSON(w, errors.New("unknown action"))
	}
}

func (cfg *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	request, err := http.NewRequest(http.MethodPost, "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		cfg.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		cfg.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusBadRequest {
		cfg.errorJSON(w, errors.New("invalid creds"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		cfg.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		cfg.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		cfg.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	cfg.writeJSON(w, http.StatusAccepted, payload)
}

func (cfg *Config) logItem(w http.ResponseWriter, log LogPayload) {
	jsonData, _ := json.MarshalIndent(log, "", "\t")

	request, err := http.NewRequest(http.MethodPost, "http://logger-service/log", bytes.NewBuffer(jsonData))
	if err != nil {
		cfg.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		cfg.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		cfg.errorJSON(w, errors.New("error calling log service"))
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged successfully"

	cfg.writeJSON(w, http.StatusAccepted, payload)
}
