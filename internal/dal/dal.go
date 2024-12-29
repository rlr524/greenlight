/*
internal/dal/dal.go
Created 12/29/24
Rob Ranf
rob@emiyaconsulting.com
- The dal.go file is a wrapper for all data access layer files and can be
thought of as a data access layer interface that will be implemented by all dal types
*/

package dal

import (
	"database/sql"
	"errors"
)

// ErrRecordNotFound defines a custom error and returns from any Get()
// method when looking up a record that doesn't exist in the database.
var (
	ErrRecordNotFound = errors.New("record not found")
)

// The DataAccessLayers struct wraps the MovieDAL and all additional data access layer types.
type DataAccessLayers struct {
	Movies MovieDAL
}

func NewDALs(db *sql.DB) DataAccessLayers {
	return DataAccessLayers{
		Movies: MovieDAL{DB: db},
	}
}
