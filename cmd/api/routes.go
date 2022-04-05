package main

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

//this function is gonna secure our route
func (app *application) wrap(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		//pass httprouter.Params to request context
		ctx := context.WithValue(r.Context(), "params", ps)
		//call next middleware with new context
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (app *application) routes() http.Handler {
	//initiate router
	router := httprouter.New()

	//adding middleware to the chain variable.We can add as many middlewares we want in here
	secure := alice.New(app.checkToken)

	//usual routing procedure.method,url,function
	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)

	router.HandlerFunc(http.MethodPost, "/v1/signin", app.SignIn)

	router.HandlerFunc(http.MethodGet, "/v1/movies", app.getAllMovies)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.getOneMovie)

	//Finally securing our route
	router.POST("/v1/admin/editmovie",app.wrap(secure.ThenFunc(app.editMovie)))
	// router.HandlerFunc(http.MethodPost, "/v1/admin/editmovie", app.editMovie)

	router.HandlerFunc(http.MethodGet, "/v1/admin/deletemovie/:id", app.deleteMovie)

	router.HandlerFunc(http.MethodGet, "/v1/genres", app.getAllGenres)
	router.HandlerFunc(http.MethodGet, "/v1/genres/:genre_id", app.getAllMoviesByGenre)
	return app.enableCORS(router)
}
