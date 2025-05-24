package main

import (
	"log"
	"net/http"
)

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	type mailMessage struct {
		From    string `json:"from"`
		TO      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	requestPayload := mailMessage{}
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		log.Println("Error reading JSON:", err)
		app.errorJSON(w, err)
		return
	}

	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.TO,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}
	err = app.Mailler.SendSMTPMessage(msg)
	if err != nil {
		log.Println("Error sending email:", err)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: "Sent to " + requestPayload.TO,
	}
	app.writeJSON(w, http.StatusAccepted, payload)
	// err = app.SendMail(msg)
}
