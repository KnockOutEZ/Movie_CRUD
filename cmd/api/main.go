package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

//software version
const version = "1.0.0"

//server configuration struct
type config struct {
	port int
	env  string
	//database struct
	db struct {
		//database connection string to connect with database
		dsn string
	}
}

//status struct for status request
type AppStatus struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

type application struct {
	config config
	logger *log.Logger
}

func main() {
	var cfg config

	//pushing data to config struct
	flag.IntVar(&cfg.port, "port", 8080, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment (development|production")
	//read connection from command flag.The format of the connection link it postgres:://username:pass@localhost or ip/ dbname
	flag.StringVar(&cfg.db.dsn, "dsn", "postgres://postgres:root@localhost/movie_api?sslmode=disable", "Postgres connection string")
	flag.Parse()

	//logs data in the terminal
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	//calling openDB function we just declared below
	db, err := openDB(cfg)
	if err != nil {
		//kills connection
		logger.Fatal(err)
	}
	//closes connection with database after everything is done.(Always have to do it after openning a connection with database)
	defer db.Close()

	//pushing data in application struct
	app := &application{
		config: cfg,
		logger: logger,
	}

	//default server setup build it http.server method
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	//logs server starting message in terminal
	logger.Println("Starting server on port", cfg.port)

	//listens to errors while starting server
	err = srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	//connection to database."postgres" is the connector here and then the dsn
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	//Many functions in Go use the context package to gather additional information about the environment they’re being executed in
	//cancel is for cancelling a context.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//PingContext verifies the connection to the database is still alive.
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil

}
