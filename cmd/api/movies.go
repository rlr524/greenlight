package main

import (
	"fmt"
	"github.com/rlr524/greenlight/internal/data"
	"net/http"
	"time"
)

// The createMovieHandler will handle POST actions to the /v1/movies endpoint.
func (app *Application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	// Anonymous struct to hold the information that is expected to be in the HTTp request body. This is
	// the target decode destination. This struct is not being used outside this handler and is not exported, however
	// the fields themselves need to be exported to be available to the encoding/json package.
	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	// Use the readJSON() helper to decode the request body into the input struct. If this returns an
	// error, send the client an error message along with a 400 status code.
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Dump the contents of the input struct in an HTTP response.
	_, e := fmt.Fprintf(w, "%+v\n", input)
	if e != nil {
		return
	}
}

// The showMovieHandler will handle GET actions to the /v1/movies endpoint using an id parameter.
func (app *Application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	// Use the readIDParam helper to parse the URL parameters.
	// If there is an error, send a 404 and return out of the function.
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
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
	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
