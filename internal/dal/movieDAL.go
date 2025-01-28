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
	"errors"
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
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, title, year, runtime, genres, version
		FROM movies
		WHERE id = $1 AND deleted NOT IN (true)`

	var movie model.Movie

	err := m.DB.QueryRow(query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &movie, nil
}

func (m MovieDAL) Update(movie *model.Movie) error {
	query := `
		UPDATE movies
		SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
		WHERE id = $5
		RETURNING version`

	args := []any{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
	}

	return m.DB.QueryRow(query, args...).Scan(&movie.Version)

	// So what's happening here? Notice the Update method (it's a method because it has a
	// receiver argument (m MovieDAL) for the MovieDAL type, so it is a method on MovieDAL) takes
	// a pointer to the Movie struct in the movie.go model file. This is how the data is
	// initially saved in memory when it's entered in the client. Because it's a pointer to the
	// movie data, it mutates it in place in memory, it doesn't make a new copy of the movie
	// data. Once that entry is complete, the query, with its arguments, is committed to the
	// database by the return statement. The QueryRow function is from the sql package, and it
	// executes a query on the database that returns a maximum of one row, which in our case
	// will be the updated version of the movie we just updated.
}

func (m MovieDAL) Delete(id int64) error {
	return nil
}
