package main

import (
	"errors"
	"net/http"
	"strconv"

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
		app.errorJSON(w, err)
		//breaks out from condition and not the function
		return
	}

	//just logs the id in terminal when you hit a url with id
	// app.logger.Println("id is", id)

	//grab data from database from get function in models
	movie, err := app.models.DB.Get(id)

	//dummy movie struct data
	// movie := models.Movie{
	// 	ID:          id,
	// 	Title:       "John Wick",
	// 	Description: "Shera movie mamo",
	// 	Year:        2022,
	// 	ReleaseDate: time.Now(),
	// 	Runtime:     100,
	// 	Rating:      5,
	// 	MPAARating:  "pegi-18",
	// 	Created_At:  time.Now(),
	// 	Updated_At:  time.Now(),
	// }

	//sends data to writeJSON func in utilities and it return json with a key
	err = app.writeJSON(w, http.StatusOK, movie, "movie")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) getAllMovies(w http.ResponseWriter, r *http.Request) {
	//getting all the movies.
	movies, err := app.models.DB.All()

	//checking for error
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	//finally pass the movies data to browser by using writeJSON function in utilities
	err = app.writeJSON(w, http.StatusOK, movies, "movies")

	//checking for error
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

//for getting all genres
func (app * application) getAllGenres(w http.ResponseWriter, r *http.Request){
	genres,err := app.models.DB.GenreAll()

	if err != nil {
		app.errorJSON(w,err)
		return
	}

	err = app.writeJSON(w,http.StatusOK,genres,"genres")
}

//get all movies sorted by genre
func (app *application) getAllMoviesByGenre(w http.ResponseWriter, r *http.Request){
	//to get the id from url
	params := httprouter.ParamsFromContext(r.Context())

	//convert url param to int
	genreID, err := strconv.Atoi(params.ByName("genre_id"))
	if err != nil{
		app.errorJSON(w,err)
		return
	}

	//finally calling All() function with genre id for getting movies with same genre 
	movies,err := app.models.DB.All(genreID)
	if err != nil{
		app.errorJSON(w, err)
		return
	}

	//then passing all the data to writeJSON func in utilities for showing them in browser
	err = app.writeJSON(w,http.StatusOK,movies,"movies")
}

func (app *application) deleteMovie(w http.ResponseWriter, r *http.Request){

}

func (app *application) insertMovie(w http.ResponseWriter, r *http.Request){
	
}

func (app *application) updateMovie(w http.ResponseWriter, r *http.Request){
	
}

func (app *application) searchMovies(w http.ResponseWriter, r *http.Request){
	
}