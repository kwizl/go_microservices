package main

import (
	"net/http"
)

func (app *Application) authenticate(w http.ResponseWriter, _ *http.Request) {
	payload := JsonResponse{
		Error:   false,
		Message: "Authentication Service is running",
	}

	err := app.writeJson(w, http.StatusOK, payload)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
