package main

import (
	"errors"
	"log"
	"net/http"
)

func (app *Config) authenticate(w http.ResponseWriter, r *http.Request) {
	log.Println("Auth method init")
	var reqPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &reqPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	// validate user against database
	user, err := app.Models.User.GetByEmail(reqPayload.Email)
	if err != nil {
		log.Println("Invalid password", reqPayload.Password, reqPayload.Email)

		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// valid
	valid, err := user.PasswordMatches(reqPayload.Password)
	if err != nil || !valid {
		log.Println("Invalid password", valid, reqPayload.Password)
		app.errorJSON(w, errors.New("invalid Creds"), http.StatusBadRequest)
		return

	}
	payload := JsonResponse{
		Error: false,
		// Message: fmt.Sprintf("Logged in user %v", user.Email),
		// Data:    user.,
	}

	app.writeJSON(w, http.StatusAccepted, payload)

}
