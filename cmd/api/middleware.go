package main

import (
	"encoding/json"
	"net/http"
)

//that how all middlewares are used to prevent cors error
func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		json.NewEncoder(w).Encode("OKOK")
		next.ServeHTTP(w, r)
	})
}
