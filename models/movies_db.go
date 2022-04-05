package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)


type DBModel struct {
	DB *sql.DB
}

//Get returns one movie and err if any from database
func (m *DBModel) Get(id int) (*Movie, error) {
	//context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//query for database. id=$1 is the placeholder
	query := `select id, title, description, year, release_date, rating, runtime, mpaa_rating,
	created_at, updated_at from movies where id = $1
`

	//using query variable to populate the row variable.
	//QueryRowContext executes a query that is expected to return at most one row. QueryRowContext always returns a non-nil value. Errors are deferred until Row's Scan method is called. If the query selects no rows, the *Row's Scan will return ErrNoRows. Otherwise, the *Row's Scan scans the first selected row and discards the rest.
	row := m.DB.QueryRowContext(ctx, query, id)

	//imagine this as the Movie struct
	var movie Movie
	//scans the input texts which is given in the standard input, reads from there and stores the successive space-separated values into successive arguments
	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.Description,
		&movie.Year,
		&movie.ReleaseDate,
		&movie.Rating,
		&movie.Runtime,
		&movie.MPAARating,
		&movie.Created_At,
		&movie.Updated_At,
	)
	if err != nil {
		return nil, err
	}

	//get the genres, if any

	//for querying one genre for row variable
	//mg is just an alias for movies_genre.and g is an alias for genres table
	query = `
	select
		mg.id,mg.movie_id,mg.genre_id,g.genre_name
	from
		movies_genres mg
		left join genres g on (g.id = mg.genre_id)
	where
		mg.movie_id = $1
	`
	//gives me a specific row depending on my id provided.
	rows, _ := m.DB.QueryContext(ctx, query, id)
	//closing the context to avoid any resource leaks.
	defer rows.Close()

	genres := make(map[int]string)

	//Next() function is used to get the next element in list go golang.
	for rows.Next() {
		var mg MovieGenre
		//scans from rows
		err := rows.Scan(
			&mg.ID,
			&mg.MovieID,
			&mg.GenreID,
			&mg.Genre.GenreName,
		)
		if err != nil {
			return nil, err
		}
		//setting our map values
		genres[mg.ID] = mg.Genre.GenreName
	}

	movie.MovieGenre = genres

	//we are returning a movie reference cause we are returning pointer.
	return &movie, nil
}

//All() returns all movies and if serched by genre it will show all movies with same genre from database
func (m *DBModel) All(genre ...int) ([]*Movie, error) {
	//setup our context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//sort by genre functionality starts from here
	where := ""
	//if there is an argument passed calling this All function this snippet will run and give us a query for getting all movies with same genre.
	//genre[0] will be whichever id we supplied while calling this function
	if len(genre) > 0 {
		where = fmt.Sprintf("where id in (select movie_id from movies_genres where genre_id = %d)", genre[0])
	}

	//query for getting all movies ordered by title.And if we search by same genre it will put the "where" variable data with query in this "query" variable.
	query := fmt.Sprintf(`select id, title, description, year, release_date, rating, runtime, mpaa_rating,
	created_at, updated_at from movies %s order by title`, where)

	//sort by genre functionality ends here

	//store that query result in the rows variable
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	//always remember to close context
	defer rows.Close()

	//movies variable will hold the final result
	var movies []*Movie

	//iterate through the rows variable
	for rows.Next() {
		var movie Movie
		//scan the currect row and push data into our movie
		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Description,
			&movie.Year,
			&movie.ReleaseDate,
			&movie.Rating,
			&movie.Runtime,
			&movie.MPAARating,
			&movie.Created_At,
			&movie.Updated_At,
		)
		if err != nil {
			return nil, err
		}

		//get the genres, if any

		//for querying one genre for row variable
		//mg is just an alias for movies_genre.and g is an alias for genres table
		genreQuery := `
	select
		mg.id,mg.movie_id,mg.genre_id,g.genre_name
	from
		movies_genres mg
		left join genres g on (g.id = mg.genre_id)
	where
		mg.movie_id = $1
	`
		//gives me a specific row depending on my id provided.
		GenreRows, _ := m.DB.QueryContext(ctx, genreQuery, movie.ID)
		//closing the context to avoid any resource leaks.
		defer rows.Close()

		//creating a map for genres
		genres := make(map[int]string)

		//Next() function is used to get the next element in list go golang.
		for GenreRows.Next() {
			var mg MovieGenre
			//scanning the GenreRows data into MovieGenre variable
			err := GenreRows.Scan(
				&mg.ID,
				&mg.MovieID,
				&mg.GenreID,
				&mg.Genre.GenreName,
			)
			if err != nil {
				return nil, err
			}
			//setting our map values
			genres[mg.ID] = mg.Genre.GenreName
		}
		//closing GenreRows when we are done looping through the GenreRows.
		GenreRows.Close()

		//assign genres to the MovieGenre field of the movie variable
		movie.MovieGenre = genres
		//append movie reference to the movies variable to our slice of movies
		movies = append(movies, &movie)
	}

	//and finally return all data
	return movies, nil
}

//for getting all genres from database
func (m *DBModel) GenreAll() ([]*Genre, error) {
	//setup our context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//query for getting all genres database.
	query := `select id, genre_name, created_at, updated_at from genres order by genre_name
`

	//using query variable to populate the row variable.
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []*Genre

	for rows.Next() {
		var g Genre
		err := rows.Scan(
			&g.ID,
			&g.GenreName,
			&g.Created_At,
			&g.Updated_At,
		)
		if err != nil {
			return nil, err
		}
		genres = append(genres, &g)
	}

	return genres, nil
}

//inserts new movie in the database
func (m *DBModel) InsertMovie(movie Movie) error {
	//setup our context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//query for adding new movies in database.
	query := `insert into movies (title, description, year, release_date, runtime, rating, mpaa_rating,
			  created_at,updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	
	//adding data from editMovie to database
	_, err := m.DB.ExecContext(ctx, query,
		movie.Title,
		movie.Description,
		movie.Year,
		movie.ReleaseDate,
		movie.Runtime,
		movie.Rating,
		movie.MPAARating,
		movie.Created_At,
		movie.Updated_At,
	)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

//for updating movies in database
func (m *DBModel) UpdateMovie(movie Movie) error {
	//setup our context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//query for updating data in database.
	query := `update movies set title = $1, description = $2, year = $3, release_date = $4, runtime = $5, rating = $6, mpaa_rating = $7,
			  updated_at = $8 where id = $9`
	
	//adding data from editMovie to database
	_, err := m.DB.ExecContext(ctx, query,
		movie.Title,
		movie.Description,
		movie.Year,
		movie.ReleaseDate,
		movie.Runtime,
		movie.Rating,
		movie.MPAARating,
		movie.Updated_At,
		movie.ID,
	)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}