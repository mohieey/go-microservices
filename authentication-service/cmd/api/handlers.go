package main

import (
	"errors"
	"log"
	"net/http"
)

func (cfg *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var reqPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := cfg.readJSON(w, r, &reqPayload)
	if err != nil {
		cfg.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	user, err := cfg.Models.User.GetByEmail(reqPayload.Email)
	if err != nil {
		cfg.errorJSON(w, errors.New("invalid creds"), http.StatusBadRequest)
		return
	}

	isValidPassword, err := user.PasswordMatches(reqPayload.Password)
	log.Println("================================================================")
	log.Println(err)
	log.Println(isValidPassword)
	log.Println("================================================================")
	if err != nil || !isValidPassword {
		cfg.errorJSON(w, errors.New("invalid creds"), http.StatusBadRequest)
		return
	}

	resPayload := jsonResponse{
		Error:   false,
		Message: "logged in successfully",
		Data:    user,
	}

	cfg.writeJSON(w, http.StatusAccepted, resPayload)
}
