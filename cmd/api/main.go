package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)
//software version
const version = "1.0.0"

//server configuration struct
type config struct {
	port int
	env  string
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
	flag.Parse()

	//logs data in the terminal
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

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
	err := srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
