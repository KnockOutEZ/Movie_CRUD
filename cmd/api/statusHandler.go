package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) statusHandler(w http.ResponseWriter, r *http.Request) {
	//dummy data for /status url
	currentStatus := AppStatus{
		Status:      "Available",
		Environment: app.config.env,
		Version:     version,
	}

	//converts currentStatus struct to json
	js, err := json.MarshalIndent(currentStatus, "", "\t")
	//prints error if occurs
	if err != nil {
		app.logger.Println(err)
	}

	//sends the json data to the browser
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
