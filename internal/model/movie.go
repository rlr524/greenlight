package model

import (
	"github.com/rlr524/greenlight/internal/validator"
	"time"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Title     string    `json:"title"`
	Year      int32     `json:"year"`
	Runtime   Runtime   `json:"runtime,omitempty"`
	Genres    []string  `json:"genres"`
	Deleted   bool      `default:"false" json:"deleted"`
	Version   int32     `json:"version"`
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	currentYear := int32(time.Now().Year())

	// Use the Check() method to execute the validation checks. This will add the provided key and
	// error message to the errors map if the check does not evaluate to true. For example, in the
	// first line, "check that the title is not equal to the empty string". In the second, "check
	// that the length of the title is less than or equal to 500 bytes" and so on.
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes (about "+
		"500 characters) long")

	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= (currentYear+2), "year", "must not be more than two "+
		"years in the future")

	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive whole number")

	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least one genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than five genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}
