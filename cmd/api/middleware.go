package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pascaldekloe/jwt"
)

//that how all middlewares are used to prevent cors error
func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//setting header for cors
		w.Header().Set("Access-Control-Allow-Origin", "*")
		//modifying header so we can allow certain things to be passed
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		next.ServeHTTP(w, r)
	})
}


//this middleware will check for valid jwt tokens
func (app *application) checkToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//adding a header.Vary is the key and Authorization is the value
		w.Header().Add("Vary", "Authorization")

		//getting Authorization from header.The header we will have to send from our frontend
		authHeader := r.Header.Get("Authorization")


		if authHeader == "" {
			// could set an anonymous user
		}

		//taking headerParts variable and splitting it into spaces.And we will get 2 parts from this split.
		headerParts := strings.Split(authHeader, " ")
		//Now we check the header
		if len(headerParts) != 2 {
			app.errorJSON(w, errors.New("invalid auth header"))
			return
		}

		if headerParts[0] != "Bearer" {
			app.errorJSON(w, errors.New("unauthorized - no bearer"))
			return
		}

		//Now we start checking the token itself
		token := headerParts[1]

		claims, err := jwt.HMACCheck([]byte(token), []byte(app.config.jwt.secret))
		if err != nil {
			// http.StatusForbidden returns 403 (Forbidden) Status Code in HTTP response 
			app.errorJSON(w, errors.New("unauthorized - failed hmac check"),http.StatusForbidden)
			return
		}

		if !claims.Valid(time.Now()) {
			app.errorJSON(w, errors.New("unauthorized - token expired"),http.StatusForbidden)
			return
		}

		if !claims.AcceptAudience("mydomain.com") {
			app.errorJSON(w, errors.New("unauthorized - invalid audience"),http.StatusForbidden)
			return
		}

		if claims.Issuer != "mydomain.com" {
			app.errorJSON(w, errors.New("unauthorized - invalid issuer"),http.StatusForbidden)
			return
		}

		userID, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			app.errorJSON(w, errors.New("unauthorized"),http.StatusForbidden)
			return
		}

		log.Println("Valid user:", userID)

		next.ServeHTTP(w, r)
	})
}
