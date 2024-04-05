package main

import (
	"fmt"
	"github.com/rlr524/greenlight/internal/model"
	"net/http"
	"time"
)

// createMovieHandler() creates a new movie.
// Method: POST
// Endpoint: /v1/movies
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintln(w, "create a new movie")
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
