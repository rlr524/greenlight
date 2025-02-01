package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	// Initialize a new httprouter instance
	r := httprouter.New()
	v := "/v1"

	r.NotFound = http.HandlerFunc(app.notFoundResponse)
	r.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowed)

	// Register the relevant methods, URL patterns and handler functions for the endpoints
	// using the HandlerFunc() method.
	r.HandlerFunc(http.MethodGet, v+"/healthcheck", app.healthcheckHandler)
	r.HandlerFunc(http.MethodPost, v+"/movies", app.createMovieHandler)
	r.HandlerFunc(http.MethodGet, v+"/movies/:id", app.getMovieHandler)
	r.HandlerFunc(http.MethodPut, v+"/movies/:id", app.updateMovieHandler)

	return app.recoverPanic(r)
}
