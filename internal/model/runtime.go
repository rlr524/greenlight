package model

import (
	"fmt"
	"strconv"
)

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
