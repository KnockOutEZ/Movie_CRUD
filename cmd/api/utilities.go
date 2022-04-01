package main

import (
	"encoding/json"
	"net/http"
)

//for converting data into json and sending them to browser
func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, wrap string) error {
	//wraps my content with a key
	wrapper := make(map[string]interface{})

	//wrap is the key
	wrapper[wrap] = data

	//converts wrapper content to json
	js, err := json.Marshal(wrapper)
	//return error if occurs
	if err != nil {
		return err
	}

	//sends wrapper json data to browser
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	//return nothing.Have to use this cause we have a return type in our func
	return nil
}

//for better error handling
func (app *application) errorJSON(w http.ResponseWriter,err error){
	//struct for error message
	type jsonError struct{
		Message string `json:"message"`
	}

	//setting error struct value
	theError := jsonError{
		//sets the error message as a string
		Message: err.Error() ,
	}

	//finally sends all the info to writeJSON function
	app.writeJSON(w, http.StatusBadRequest,theError,"error")
}