package main

import (
	"fmt"
	"net/http"
)

// The logError() method is a generic helper for logging an error message along
// with the current request method and URL as attributes in the log entry.
func (app *application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

// The errorResponse() method is a generic helper for sending JSON-formatted error messages to the
// client with a given status code. Note the *any* type for the message parameter, rather than
// a string; this provides more flexibility over the values that can be contained in the response.
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	// Write the response using the writeJSON() helper. If this happens to return an error then
	// log it and fall back to sending the client an empty response with a 500 status code.
	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

// The serverErrorResponse() method is used when the application encounters an unexpected problem
// at runtime. It logs the detailed error message, then uses the errorResponse() helper to send
// a 500 error and JSON response (containing a generic error message) to the client.
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// The notFoundResponse() method will be used to send a 404 status and JSON response to the client.
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// The methodNotAllowed() method will be used to send a 405 status and JSON response to the client.
func (app *application) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	var method = r.Method
	message := fmt.Sprintf("the %s method is not supported for this resource", method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

// The badRequestResponse method will be used to send a 400 status and JSON response to the client.
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

// The failedValidationResponse() method  will be used to write the 422 Unprocessable Entity response
// and the contents of the errors map from the Validator type as JSON response body. Note that the
// errors parameter has the type map[string]string which is exactly the same as the errors map
// in the Validator type.
func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request,
	errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}
