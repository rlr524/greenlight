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

type envelope map[string]any

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int,
	data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		app.logger.Error(err.Error())
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(js)
	if err != nil {
		app.logger.Error(err.Error())
	}

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	// Allow requests of a maximum of only 1MB
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Decode the request body into the target destination. Initialize a new
	// json.Decoder instance which reads from the request body, and then use the
	// Decode() method to decode the body contents into the input struct. Importantly, note that
	// when Decode() is called, a *pointer to the input struct is passed as the
	// target decode destination. If there was an error during decoding, the generic
	// errorResponse() helper is used to send the client a 400 Bad Request response
	// containing the error message. When calling Decode(), you must pass a non-nil
	// pointer as the target decode destination. If you don't use a pointer, it will
	// return a json.InvalidUnmarshalError error at runtime.
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(dst)
	if err != nil {
		// If there is an error during decoding, start triaging the errors...
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		// Use the errors.As() function to check whether the error has the type
		// *json.SyntaxError. If it does, then return a plain-English error message which
		// includes the location of the problem.
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly formed JSON as character %d",
				syntaxError.Offset)
		// In some circumstances Decode() may also return an io.ErrUnexpectedEOF error for syntax
		// errors in the JSON. So we check for this using errors.Is() and return a generic error
		// message. (Open issue regarding this at https://github.com/golang/go/issues/25956)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly formed JSON")

		// Likewise, catch any *json.UnmarshalTypeError errors. These occur when the JSON
		// value is the wrong type for the target destination. If the error relates to a
		// specific field, then we include that in our error message to make it easier
		// for the client to debug.
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				// The %q verb when used for a value that returns a string is a double-quoted
				// string, opposed to %s which is the uninterpreted bytes of a string or a slice.
				return fmt.Errorf("body contains incorrect JSON type for field %q",
					unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type at character %d",
				unmarshalTypeError.Offset)

		// An io.EOF error will be returned by Decode() if the request body is empty. We
		// check for this with errors.Is() and return a plain-english error message
		// instead.
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// If the json contains a field which cannot be mapped to the target destination, then
		// Decode() will return an error message in the format "json: unknown field "<name>". We
		// check for this, extract the field name from the error, and interpolate it into a custom
		// error message. Note that there's an open issue at https://github.com/golang/go/issues/29035
		// regarding turning this into a distinct error type in the future.
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		// Use the errors.As() function to check whether the error has the type *http.MaxBytesError.
		// If it does, then the request body exceeded the size limit of 1MB and a clear
		// error message is returned.
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)

		// A json.InvalidUnmarshalError error will be returned if we pass something that
		// is not a non-nil pointer to Decode(). We catch this and panic, rather than
		// returning an error to our handler. At the end of this chapter we'll talk about
		// panicking versus returning errors, and discuss why it's an appropriate thing
		// to do in this specific situation.
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	// Call Decode() again, using a pointer to an empty anonymous struct as the destination. if the
	// request body only contained a single JSON value this will return an io.EOF error. If anything
	// else is returned, then there is additional data in the request body and a custom
	// error message is returned.
	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}
