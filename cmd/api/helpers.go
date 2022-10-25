package main

import (
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// Retrieve the "id" URL parameter from the current request context, then convert it to an integer and return
// it. If the operation isn't successful, return 0 and an error
func (app *Application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func (app *Application) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	// Encode the data to JSON, returning an error if there was one
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Append a newline for terminal reading
	js = append(js, '\n')

	// Loop through header map and add headers to the http.ResponseWriter header map. Note that it's OK
	// if the provided header map is nil. Go doesn't throw an error if you try to range over
	// (or generally, read from) a nil map.
	for key, value := range headers {
		w.Header()[key] = value
	}

	// Add the content-type header then write the status code and JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// The Let's Go book does not handle Write errors since at this point we should know there won't be any
	// errors since we're doing error checking at the point of marshaling the data, however, it's still best
	// to handle any errors that may happen between the data encoding and the write function.
	_, e := w.Write(js)
	if e != nil {
		panic(e)
	}

	return nil
}
