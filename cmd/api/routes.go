package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {

	r := httprouter.New()

	r.NotFound = http.HandlerFunc(app.notFoundResponse)
	r.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	r.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	r.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)
	r.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)

	return r
}
