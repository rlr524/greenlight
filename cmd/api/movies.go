package main

import (
	"fmt"
	"net/http"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintln(w, "create a new movie")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	
	_, err = fmt.Fprintf(w, "show the details of movie %d\n", id)
	if err != nil {
		app.logger.Info("Error: %v\n", err)
		return
	}
}
