package main

import (
	"backend/models"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

func (app *application) getOneMovie(w http.ResponseWriter, r *http.Request) {

	//default julienschmidt/httprouter parameter setup.It gets key from url
	params := httprouter.ParamsFromContext(r.Context())

	//converts the id from string to int
	id, err := strconv.Atoi(params.ByName("id"))
	//checks if the parameter id passed is valid
	if err != nil {
		//prints parameter error in terminal
		app.logger.Print(errors.New("invalid id parameter"))
		//sends param to errorJSON func in utilities for better error handling 
		app.errorJSON(w,err)
		//breaks out from condition and not the function
		return
	}

	//just logs the id in terminal when you hit a url with id
	app.logger.Println("id is", id)

	//dummy movie struct data
	movie := models.Movie{
		ID:          id,
		Title:       "John Wick",
		Description: "Shera movie mamo",
		Year:        2022,
		ReleaseDate: time.Now(),
		Runtime:     100,
		Rating:      5,
		MPAARating:  "5",
		Created_At:  time.Now(),
		Updated_At:  time.Now(),
	}

	//sends data to writeJSON func in utilities and it return json with a key
	err = app.writeJSON(w, http.StatusOK,movie,"movie")
}

func (app *application) getAllMovies(w http.ResponseWriter, r *http.Request) {

}
