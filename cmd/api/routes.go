package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	//initiate router
	router := httprouter.New()

	//usual routing procedure.method,url,function
	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)

	router.HandlerFunc(http.MethodGet, "/v1/movie", app.getAllMovies)
	router.HandlerFunc(http.MethodGet, "/v1/movie/:id", app.getOneMovie)

	return router
}
