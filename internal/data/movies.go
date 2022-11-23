package data

import (
	"time"

	"github.com/BuddhiLW/greenlight/internal/validator"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year, omitempty"`
	Runtime   Runtime   `json:"runtime, omitempty"`
	Genres    []string  `json:"genres, omitempty"`
	Version   int32     `json:"version"`
}

func ValidateMovie(v *validator.Validator, movie *Movie) {

	// Check if title is not empty. If it is empty (check is False), add title-type error
	v.Check(!(movie.Title == ""), "title", "Must be provided")
	// Check if title is not less then 500 bytes. If it is greater (check is False), add title-type error
	v.Check(!(len(movie.Title) <= 500), "title", "Must not be more than 500 bytes long")

	// Check if runtime is a positive, non-zero, number
	v.Check(!(movie.Runtime < 0), "runtime", "Must be a positive integer")
	v.Check(!(movie.Runtime == 0), "runtime", "Must be provided")

	// Check if year value is not empty
	v.Check(!(movie.Year == 0), "year", "Must be provided")
	// Check if year is not lesser than 1888
	v.Check(!(movie.Year < 1888), "year", "Must be greater than 1888")
	// Check if movie-year is not a value greater than current year it's being added
	v.Check(!(movie.Year >= int32(time.Now().Year())), "year", "Must not be in the future")

	// Check if Genre is empty
	v.Check(!(movie.Genres == nil), "genres", "Must be provided")
	// Check if there is at least one Genre
	v.Check(!(len(movie.Genres) <= 1), "genres", "Must contain at lest one genre")
	// Check if there is less than five Genres
	v.Check(!(len(movie.Genres) >= 5), "genres", "Must not contain more than 5 genres")
	// Check the uniqueness of each Genre
	v.Check(validator.Unique(movie.Genres), "genres", "Must not contain duplicates")
}
