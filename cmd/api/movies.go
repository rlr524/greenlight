package main

import (
	"fmt"
	"github.com/rlr524/greenlight/internal/data"
	"net/http"
	"time"
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

	// A new instance of the Movie struct, containing the ID that was extracted from the URL.
	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	// Encode the movie struct instance to JSON and send it to the HTTP response
	err = app.writeJSON(w, http.StatusOK, movie, nil)
	if err != nil {
		app.logger.Print(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
