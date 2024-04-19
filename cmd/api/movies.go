package main

import (
	"encoding/json"
	"fmt"
	"github.com/rlr524/greenlight/internal/model"
	"net/http"
	"time"
)

// createMovieHandler() creates a new movie.
// Method: POST
// Endpoint: /v1/movies
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	// Declare an anonymous struct to hold the information that is expected to be in the
	// HTTP request body. This truct is the *target decode destination.
	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	// Initialize a new json.Decoder instance which reads from the request body, and then use the
	// Decode() method to decode the body contents into the input struct. Importantly, note that
	// when Decode() is called, a *pointer to the input struct is passed as the
	// target decode destination. If there was an error during decoding, the generic
	// errorResponse() helper is used to send the client a 400 Bad Request response containing
	// the error message. When calling Decode(), you must pass a non-nil pointer as the target
	// decode destination. If you don't use a pointer, it will return a json.InvalidUnmarshalError
	// error at runtime.
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// Dump the contents of the input struct in an HTTP response.
	_, _ = fmt.Fprintf(w, "%+v\n", input)
}

// showMovieHandler() retrieves the details of a specific movie by its ID.
// Method: GET
// Endpoint: /v1/movies/:id
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie := model.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.logger.Error(err.Error())
		app.serverErrorResponse(w, r, err)
	}
}
