package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/BuddhiLW/greenlight/internal/data"
	"github.com/BuddhiLW/greenlight/internal/validator"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	v := validator.New()
	// Check if title is not empty. If it is empty (check is False), add title-type error
	v.Check(!(input.Title == ""), "title", "Must be provided")
	// Check if title is not less then 500 bytes. If it is greater (check is False), add title-type error
	v.Check(!(len(input.Title) <= 500), "title", "Must not be more than 500 bytes long")

	v.Check(!(input.Runtime < 0), "runtime", "Must be a positive integer")
	v.Check(!(input.Runtime == 0), "runtime", "Must be provided")

	// Check if year value is not empty
	v.Check(!(input.Year == 0), "year", "Must be provided")
	// Check if year is not lesser than 1888
	v.Check(!(input.Year < 1888), "year", "Must be greater than 1888")
	// Check if movie-year is not a value greater than current year it's being added
	v.Check(!(input.Year >= int32(time.Now().Year())), "year", "Must not be in the future")

	v.Check(!(input.Genres == nil), "genres", "Must be provided")
	v.Check(!(len(input.Genres) <= 1), "genres", "Must contain at lest one genre")
	v.Check(!(len(input.Genres) >= 5), "genres", "Must not contain more than 5 genres")
	v.Check(validator.Unique(input.Genres), "genres", "Must not contain duplicates")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	envel := envelope{"movie": movie}

	err = app.writeJSON(w, http.StatusOK, envel, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
