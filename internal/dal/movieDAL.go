/*
internal/dal/movieDAL.go
Created 12/29/24
Rob Ranf
rob@emiyaconsulting.com
- The movieDAL.go file is the data access layer for the Movie
type and implements all database CRUD operations for that type.
*/

package dal

import (
	"database/sql"
	"github.com/lib/pq"
	"github.com/rlr524/greenlight/internal/model"
)

type MovieDAL struct {
	DB *sql.DB
}

// The Insert method accepts a pointer to a movie
// struct which should contain the data for the new record.
func (m MovieDAL) Insert(movie *model.Movie) error {
	// Define the SQL query for inserting a new record into the
	// movies table and returning the system-generated data.
	query := `
		INSERT INTO movies (title, year, runtime, genres) 
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`

	// Create an args slice containing the values for the placeholder
	// parameters from the movie struct. Declaring this slice immediately next to
	// the SQL query helps to make it clear what values are being used where in the query.
	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	// Use the QueryRow() method that is available from the Go database/sql library to execute
	// the SQL query on the connection pool, passing in the args slice as a variadic parameter
	// and scanning the system-generated id, created_at and version values into the movie struct.
	return m.DB.QueryRow(query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (m MovieDAL) Get(id int64) (*model.Movie, error) {
	return nil, nil
}

func (m MovieDAL) Update(movie *model.Movie) error {
	return nil
}

func (m MovieDAL) Delete(id int64) error {
	return nil
}
