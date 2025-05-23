package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	Name string `json:"name,omitempty"`
	Data string `json:"data,omitempty"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}
	_ = app.writeJSON(w, http.StatusOK, payload)

}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var reqPayload RequestPayload
	err := app.readJSON(w, r, &reqPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	log.Println("reqPayload ", reqPayload.Action)
	switch reqPayload.Action {
	case "auth":
		app.authenticate(w, reqPayload.Auth)
	case "log":
		app.logItem(w, reqPayload.Log)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) logItem(w http.ResponseWriter, l LogPayload) {
	jsonData, _ := json.MarshalIndent(l, "", "\t")
	logServiceURL := "http://logger-service/log"
	req, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling log service"))
		return
	}

	var jsonResponse jsonResponse
	jsonResponse.Error = false
	jsonResponse.Message = "log entry created successfully"

	app.writeJSON(w, http.StatusAccepted, jsonResponse)

}
func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create json and send to auth microservice
	log.Println("auth initiated ")
	jsonData, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	/// call service
	req, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("err at post ", err)
		app.errorJSON(w, err)
		return
	}
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Println("err at Do ", err)

		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// make sure we get correct status code
	fmt.Println("err at status code ", response.StatusCode)
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid creds"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling auth service"))
		return
	}
	// create var who read response body
	jsonFormService := jsonResponse{}
	// decode the json from the auth service

	err = json.NewDecoder(response.Body).Decode(&jsonFormService)
	fmt.Println("response body ", err)
	if err != nil {
		log.Println("err at decode ", err)
		app.errorJSON(w, err)
		return
	}
	if jsonFormService.Error {
		log.Println("81")
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "authenticated"
	payload.Data = jsonFormService.Data
	app.writeJSON(w, http.StatusAccepted, payload)

}
