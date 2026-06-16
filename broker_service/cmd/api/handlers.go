package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Application) broker(w http.ResponseWriter, r *http.Request) {
	payload := JsonResponse{
		Error:   false,
		Message: "Broker Service is running",
	}

	err := app.writeJson(w, http.StatusOK, payload)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *Application) handleSubmission(w http.ResponseWriter, r *http.Request) {
	var RequestPayload RequestPayload

	err := app.readJson(w, r, &RequestPayload)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	switch RequestPayload.Action {
	case "auth":
		app.authenticate(w, &RequestPayload.Auth)
	default:
		app.errorJson(w, errors.New("Invalid action"))
	}
}

func (app *Application) authenticate(w http.ResponseWriter, p *AuthPayload) {
	jsonData, _ := json.MarshalIndent(p, "", "\t")

	request, err := http.NewRequest("POST", "http://authentication:8080/autheticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err)
		return
	}
	
	client := &http.Client {Timeout: 7 * time.Second,}
	response, err := client.Do(request)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJson(w, errors.New("Invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJson(w, errors.New("Error calling authetication service"))
		return
	}

	var serviceResponse JsonResponse

	err = json.NewDecoder(response.Body).Decode(&response)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	if serviceResponse.Error {
		app.errorJson(w, err, http.StatusUnauthorized)
		return
	}

	var payload JsonResponse
	payload.Error = false
	payload.Message = "Authorized"
	payload.Data = serviceResponse

	err = app.writeJson(w, http.StatusOK, payload)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
