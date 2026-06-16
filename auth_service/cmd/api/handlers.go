package main

import (
	"errors"
	"fmt"
	"net/http"
)

func (app *Application) authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email string    `json:"email"`
		Password string `json:"password"`
	}

	// Validates the payload
	err := app.readJson(w, r, &requestPayload)
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	// Checks if user exists
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJson(w, errors.New("Invalid credentials"), http.StatusUnauthorized)
		return
	}
	
	// Checks if password matches user
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || valid {
		app.errorJson(w, errors.New("Invalid credentials"), http.StatusUnauthorized)
		return
	}

	payload := JsonResponse {
		Error: false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data: user,
	}

	err = app.writeJson(w, http.StatusAccepted, payload)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
