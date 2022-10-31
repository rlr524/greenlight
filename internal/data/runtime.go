package data

import (
	"fmt"
	"strconv"
)

// Runtime is a custom type, which has an underlying type int32 (just as the Movie struct).
type Runtime int32

// MarshalJSON is a method on the Runtime type so that it satisfies the json.Marshaler
// interface. This should return the JSON-encoded value for the movie runtime, in this case
// a string in the format "<runtime> mins". This method is intentionally using a value
// receiver, not a pointer receiver, so it can be invoked on pointers and values, where
// pointer receivers can only be invoked on pointers. Remember that Go does not have classes,
// so this is a *method* because it's called on the Runtime type using the receiver, whereas
// a function has no receiver.
func (r Runtime) MarshalJSON() ([]byte, error) {
	// Generate a string containing the movie runtime in the required format.
	jsonValue := fmt.Sprintf("%d mins", r)

	// Use the strconv.Quote() function on the string to wrap it in double quotes. It needs
	// to be surrounded by double quotes in order to be a valid JSON string.
	quotedJSONValue := strconv.Quote(jsonValue)

	// Convert the quoted string value to a byte slice and return it.
	return []byte(quotedJSONValue), nil

}
