package main

import (
	"fmt"
	"github.com/rlr524/greenlight/internal/model"
	"github.com/rlr524/greenlight/internal/validator"
	"net/http"
	"time"
)

// createMovieHandler() creates a new movie.
// Method: POST
// Endpoint: /v1/movies
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	// Declare an anonymous struct to hold the information that is expected to be in the
	// HTTP request body. This struct is the *target decode destination.
	var input struct {
		Title   string        `json:"title"`
		Year    int32         `json:"year"`
		Runtime model.Runtime `json:"runtime"`
		Genres  []string      `json:"genres"`
	}

	// Use the readJSON() helper to decode the request body into the input struct. If this
	// returns an error, send the client the error message along with a 400 status code.
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movie := &model.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	// Initialize a new Validator instance
	v := validator.New()

	// Call the ValidateMovie() function and return a response containing the errors
	// if any of the checks fail.
	if model.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
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
		Year:      1942,
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Deleted:   false,
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.logger.Error(err.Error())
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getMoviesHandler(w http.ResponseWriter, r *http.Request) {
	// Some code to come to fetch all movies that aren't status deleted
}
