package models

import (
	"database/sql"
	"time"
)

//Models is the wrapper for database
type Models struct {
	DB DBModel
}

//NewModels returns models with db pool
func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

//type for movie
type Movie struct {
	ID          int          `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Year        int          `json:"year"`
	ReleaseDate time.Time    `json:"release_date"`
	Runtime     int          `json:"runtime"`
	Rating      int          `json:"rating"`
	MPAARating  string       `json:"mpaa_rating"`
	Created_At  time.Time    `json:"-"`
	Updated_At  time.Time    `json:"-"`
	MovieGenre  map[int]string `json:"genres"`
}

//type for genre
type Genre struct {
	ID         int       `json:"-"`
	GenreName  string    `json:"genre_name"`
	Created_At time.Time `json:"-"`
	Updated_At time.Time `json:"-"`
}

//type for movie genre
// the "-" in json means they wont be delivered to frontend
type MovieGenre struct {
	ID         int       `json:"-"`
	MovieID    int       `json:"-"`
	GenreID    int       `json:"-"`
	Genre      Genre     `json:"genre"`
	Created_At time.Time `json:"-"`
	Updated_At time.Time `json:"-"`
}


//User is the type for users
type User struct{
	ID int
	Email string
	Password string
}