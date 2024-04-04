package main

import (
	"fmt"
	"net/http"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	// Note we are using a lambda function in our HandlerFunc
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event of a panic
		// as Go unwinds the stack).
		defer func() {
			// Use the builtin recover function to check if there has been a panic or not.
			if err := recover(); err != nil {
				// If there was a panic, set a "Connection: close" header on the response. This
				// acts as a trigger to make Go's http server automatically close the current
				// connection after a response has been sent.
				w.Header().Set("Connection", "close")
				// The value returned by recover() has the type *any*, so we use fmt.Errorf() to
				// normalize it into an error and call the serverErrorResponse() helper. In turn,
				// this will log the error using the custom Logger type at the ERROR level and
				// send the client a 500 response.
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
