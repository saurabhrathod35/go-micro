package main

import (
	"logger-service/data"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name,omitempty"`
	Data string `json:"data,omitempty"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var payload JSONPayload
	_ = app.readJSON(w, r, &payload)

	event := data.LogEntry{
		Name: payload.Name,
		Data: payload.Data,
	}
	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	response := jsonResponse{
		Error:   false,
		Message: "Log entry created successfully",
	}
	app.writeJSON(w, http.StatusCreated, response)

}
