package main

import (
	"errors"
	"fmt"
	"github.com/rlr524/greenlight/internal/dal"
	"github.com/rlr524/greenlight/internal/model"
	"github.com/rlr524/greenlight/internal/validator"
	"net/http"
	"strconv"
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
	err = app.dataAccessLayers.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// getMovieHandler() retrieves the details of a specific movie by its ID.
// Method: GET
// Endpoint: /v1/movies/:id
func (app *application) getMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.dataAccessLayers.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, dal.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.logger.Error(err.Error())
		app.serverErrorResponse(w, r, err)
	}
}

// updateMovieHandler updates a single movie in place
// Method: PUT
// Endpoint: /v1/movies/:id
func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.dataAccessLayers.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, dal.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// If the request contains an explicit X-Expected-Version header, verify that the movie version
	// in the database matches the expected version specified in the header.
	if r.Header.Get("X-Expected-Version") != "" {
		if strconv.Itoa(int(movie.Version)) != r.Header.Get("X-Expected-Version") {
			app.editConflictResponse(w, r)
			return
		}
	}

	var input struct {
		Title   *string        `json:"title"`
		Year    *int32         `json:"year"`
		Runtime *model.Runtime `json:"runtime"`
		Genres  []string       `json:"genres"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Copy the values from the request body to the corresponding fields of the movie record,
	// dereferencing the values for title, year, and runtime as we're passing those into the input
	// struct as pointers so we can only update a movie record if the new input value is not nil,
	// which allows for partial updates.
	if input.Title != nil {
		movie.Title = *input.Title
	}
	if input.Year != nil {
		movie.Year = *input.Year
	}
	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}
	if input.Genres != nil {
		movie.Genres = input.Genres
	}

	// Run the validator helper on the movie record.
	v := validator.New()

	if model.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Pass the updated movie record to the DAL Update method.
	err = app.dataAccessLayers.Movies.Update(movie)
	if err != nil {
		switch {
		case errors.Is(err, dal.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Write the updated movie record to return in a JSON response.
	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// getMoviesHandler fetches all movies that are not flagged as deleted
// Method: GET
// Endpoint: /v1/movies
func (app *application) getMoviesHandler(w http.ResponseWriter, r *http.Request) {
	movie, err := app.dataAccessLayers.Movies.GetAll()
	if err != nil {
		switch {
		case errors.Is(err, dal.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.logger.Error(err.Error())
		app.serverErrorResponse(w, r, err)
	}
}

// deleteMovieHandler updates the deleted flag on a single movie to true
// Method: DELETE
// Endpoint: /v1/movies/:id
func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.dataAccessLayers.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, dal.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "movie successfully deleted"},
		nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
