package main

import (
	"fmt"
	"net/http"
)

func (cfg *Config) SendEmail(w http.ResponseWriter, r *http.Request) {
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage

	err := cfg.readJSON(w, r, &requestPayload)
	if err != nil {
		cfg.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Text:    requestPayload.Message,
	}

	err = cfg.Mailer.Send(msg)
	if err != nil {
		cfg.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	cfg.writeJSON(w, http.StatusAccepted, jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("message sent successfully to %s", msg.To),
	})
}
