package main

import (
	"fmt"
	"net/http"
)

// The createMovieHandler will handle POST actions to the /v1/movies endpoint.
func (app *Application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create a new movie...")
}

// The showMovieHandler will handle GET actions to the /v1/movies endpoint using an id parameter.
func (app *Application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	// Use the readIDParam helper to parse the URL parameters.
	// If there is an error, send a 404 and return out of the function.
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// Otherwise, interpolate the movie ID in a placeholder response.
	fmt.Fprintf(w, "Show the details of movie %d\n", id)
}
