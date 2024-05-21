package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ErrInvalidRuntimeFormat is an error that the UnmarshalJSON() method can return if it's unable
// to parse or convert the JSON string successfully.
var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

type Runtime int32

// MarshalJSON here is essentially an override of the MarshalJSON function from the json library
// that is used to encode our JSON response in the writeJSON() helper method. (Remember, though
// we call json.MarshalIndent() in that method, MarshalIndent() runs json.MarshalJSON under the
// hood and then performs all the indenting.) Because this is a method on the Runtime type,
// wherever we use the Runtime type, it will use this override version of MarshalJSON in place
// of the library method.
func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", r)
	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil
}

// UnmarshalJSON is implemented on the Runtime type that satisfies the json.Unmarshaler interface.
// Because UnmarshalJSON() needs to modify the receiver (the Runtime type), it must use a pointer
// receiver, otherwise a copy of the receiver is modified, which would then be discarded
// when the method returns.
func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	// The incoming JSON value is expected to be a string in the format "<runtime> mins", so
	// the surrounding double quotes need to be removed. If this can't be done, return the
	// ErrInvalidRuntimeFormat error.
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	// Split the string to isolate the part containing the number.
	parts := strings.Split(unquotedJSONValue, " ")

	// Sanity check the parts of the string to make sure it was in the expected format. If it
	// isn't, return the ErrInvalidRuntimeFormat error again.
	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}

	// Otherwise, parse the string containing the number into an int32. Again, if this fails,
	// return the ErrInvalidRuntimeFormat error.
	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	// Convert the int32 to a Runtime type and assign this to the receiver. Note, use the *
	// operator to dereference the receiver (which is a pointer to the Runtime type) in order
	// to set the underlying value of the pointer.
	*r = Runtime(i)

	return nil
}
