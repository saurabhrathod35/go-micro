package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

	// log authentication
	err = app.logRequest("authentication", fmt.Sprintf("User %s logged in", user.Email))
	if err != nil {
		log.Println("Error logging authentication:", err)
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	app.writeJSON(w, http.StatusAccepted, payload)

}

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}
	entry.Name = name
	entry.Data = data
	jsonData, err := json.Marshal(entry)
	if err != nil {
		log.Println("Error marshalling JSON", err)
		return err
	}
	logServiceURL := "http://logger-service/logs"
	req, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error creating request:", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		log.Println("Error response from logger service:", resp.Status)
		return err

	} else {
		log.Println("Log entry sent successfully")
	}
	return nil

}
