package models

import (
	"context"
	"database/sql"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

//Get returns one movie and err if any
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
		//if genres does not exist.We append the data in genres.
		genres[mg.ID] = mg.Genre.GenreName
	}

	movie.MovieGenre = genres

	//we are returning a movie reference cause we are returning pointer.
	return &movie, nil
}

//Get returns all movies and err if any
func (m *DBModel) All(id int) ([]*Movie, error) {
	return nil, nil
}
