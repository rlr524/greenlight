package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"strconv"
	"strings"
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

// Envelope for JSON data
type envelope map[string]any

func (app *Application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	// Encode the data to JSON, returning an error if there was one. Use the MarshalIndent function
	// so that whitespace is added to the encoded JSON to make reading easier in plain interfaces such
	// as CURL.
	js, err := json.MarshalIndent(data, "", "\t")
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

func (app *Application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	// Use http.MaxBytesReader() to limit the size of the request body to 1 MB
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Init the json.Decoder and call the DisallowUnknownFields() method on it before decoding. This
	// means that if the JSON from the client now includes any field which cannot be mapped to the
	// target destination, the decoder will return an error instead of ignorning the field.
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		// If there is an error during decoding, start triage
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		// Use the errors.As() function to check whether the error has the type *json.SyntaxError.
		// If it does, then return a plain-english error message which includes the location of the problem.
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly formatted JSON (at character %d)", syntaxError.Offset)
		// In some circumstances Decode() may return an io.ErrUnexpectedEOF error for syntax errors in the
		// JSON. So, check for this using errors.Is() and return a generic error message.
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly formatted JSON")
		// Likewise, catch any *json.UnmarshalType errors. These occur when the json value is the wrong
		// type for the target destination. If the error relates to a specific field, then include it
		// in the error message to make it easier for the client to debug.
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		// An io.EOF error will be returned by Decode() if the request body is empty. Check for this with
		// errors.Is() and return a plain-english error message instead.
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		// If the JSON contains a field which cannot be mapped to the target destination, then Decode() will
		// return an error message in the format "json: unknown field "<name>"". A check is made for this, the field
		// name is extracted from the error and interpolated into the custom error message.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		// Use the errorsAs() function to check whether the error has the type *http.MaxBytesError. If it does,
		// then the request body exceeded the size limit of 1 MB and an error message is returned.
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be longer than %d bytes", maxBytesError.Limit)
		// A json.InvalidUnmarshalError error will be returned if something is passed that is a non-nil pointer
		// to Decode(). Catch this and panic, rather than returning an error to the handler.
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		// For everything else, return the error message as is.
		default:
			return err
		}
	}
	// Call Decode() again, using a pointer to any empty anonymous struct as the destination. If the request body
	// only contained a single JSON value, this will return an io.EOF error. So, if anything else is received, there
	// is additional data in the request body and a custom error message is returned.
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must contain a single JSON value")
	}

	return nil
}
