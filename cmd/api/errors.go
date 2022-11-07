package main

import (
	"fmt"
	"net/http"
)

// logError is a generic helper for logging an error message.
func (app *Application) logError(r *http.Request, err error) {
	app.logger.Print(err)
}

// errorResponse is a generic helper for sending JSON-formatted error messages to the client with a given status
// code. The *any* type is being used for the message param, rather than string, as this provides
// more flexibility over the values that can be included in the response.
func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	envelope := envelope{"error": message}

	// Write the response using the writeJSON() helper. If this happens to return an error, then log it
	// and fall back to sending the client an empty response with a 500 status code.
	err := app.writeJSON(w, status, envelope, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

// serverErrorResponse is used when the application encounters an unexpected problem at runtime. It logs the detailed
// error message, then uses the errorResponse() helper to send a 500 status code and JSON response (containing
// the generic error message) to the client.
func (app *Application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// notFoundResponse is used to send a 404 status code and JSON response to the client.
func (app *Application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"

	app.errorResponse(w, r, http.StatusNotFound, message)
}

// methodNotAllowed is used to send a 405 status code and JSON response to the client.
func (app *Application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)

	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *Application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}
