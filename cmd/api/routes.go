package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	//initiate router
	router := httprouter.New()

	//usual routing procedure.method,url,function
	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)

	router.HandlerFunc(http.MethodGet, "/v1/movies", app.getAllMovies)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.getOneMovie)

	router.HandlerFunc(http.MethodGet, "/v1/genres", app.getAllGenres)
	router.HandlerFunc(http.MethodGet, "/v1/genres/:genre_id", app.getAllMoviesByGenre)
	return app.enableCORS(router)
}
