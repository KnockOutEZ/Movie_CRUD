package main

import (
	"backend/models"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

//jsonResp is used for getting request or response status
type jsonResp struct {
	OK      bool   `json:"ok"`
	Message string `json:message`
}

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
func (app *application) getAllGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := app.models.DB.GenreAll()

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, genres, "genres")
}

//get all movies sorted by genre
func (app *application) getAllMoviesByGenre(w http.ResponseWriter, r *http.Request) {
	//to get the id from url
	params := httprouter.ParamsFromContext(r.Context())

	//convert url param to int
	genreID, err := strconv.Atoi(params.ByName("genre_id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	//finally calling All() function with genre id for getting movies with same genre
	movies, err := app.models.DB.All(genreID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	//then passing all the data to writeJSON func in utilities for showing them in browser
	err = app.writeJSON(w, http.StatusOK, movies, "movies")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) deleteMovie(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil{
		app.errorJSON(w,err)
		return
	}

	err =  app.models.DB.DeleteMovieDb(id)
	if err != nil{
		app.errorJSON(w,err)
		return
	}

	ok := jsonResp{
		OK:true,
	}

	err = app.writeJSON(w,http.StatusOK,ok,"response")
	if err != nil{
		app.errorJSON(w,err)
		return
	}
}

func (app *application) insertMovie(w http.ResponseWriter, r *http.Request) {

}

//creating a temporary payload to store all the string type data from frontend
type MoviePayload struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Year        string `json:"year"`
	ReleaseDate string `json:"release_date"`
	Runtime     string `json:"runtime"`
	Rating      string `json:"rating"`
	MPAARating  string `json:"mpaa_rating"`
}

//for adding a movie data or update existing one
func (app *application) editMovie(w http.ResponseWriter, r *http.Request) {
	var payload MoviePayload

	//taking data from request body and pushing them to payload struct temporarily (cause the data coming from frontend is all string type and we dont wanna take the hassle to convert all one by one)
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	//this is the main game.
	var movie models.Movie

	//if we edit an existing movie this code snippet will run.
	//Note that we are ignoring err by using "_"(blank variables).But we should not do that in production
	if payload.ID != "0"{
		id, _ := strconv.Atoi(payload.ID)
		m, _ := app.models.DB.Get(id)
		movie := *m
		movie.Updated_At = time.Now()
	}

	//converting and pushing all data into our main Movie struct
	movie.ID, _ = strconv.Atoi(payload.ID)
	movie.Title = payload.Title
	movie.Description = payload.Description
	movie.ReleaseDate, _ = time.Parse("2006-01-02", payload.ReleaseDate)
	movie.Year = movie.ReleaseDate.Year()
	movie.Runtime, _ = strconv.Atoi(payload.Runtime)
	movie.Rating, _ = strconv.Atoi(payload.Rating)
	movie.MPAARating = payload.MPAARating
	movie.Created_At = time.Now()
	movie.Updated_At = time.Now()

	//finally passing down the data to database
	//if we are adding a new movie in the database the "if" clause will run.And if we edit an existing one the else clause will run
	if movie.ID == 0{
		err = app.models.DB.InsertMovie(movie)
		if err != nil {
			app.errorJSON(w, err)
			return
		}
	} else{
		err = app.models.DB.UpdateMovie(movie)
		if err != nil {
			app.errorJSON(w, err)
			return
		}
	}

	//sets a response that everything worked well
	ok := jsonResp{
		OK: true,
	}

	//passing the response to browser
	err = app.writeJSON(w, http.StatusOK, ok, "response")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) searchMovies(w http.ResponseWriter, r *http.Request) {

}
