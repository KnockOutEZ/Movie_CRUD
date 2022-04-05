package main

import (
	"backend/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/pascaldekloe/jwt"
	"golang.org/x/crypto/bcrypt"
)

//creating a dummy user
var validUser = models.User{
	ID:10,
	Email:"me@gmail.com",
	Password:"$2a$14$ajq8Q7fbtFRQvXpdCq7Jcuy.Rx1h/L4J60Otx.gyNLbAYctGMJ9tK",
}
//credentials used
type Credentials struct{
	Username string `json:email`
	Password string `json:password`
}

//func for signing in
func (app *application) SignIn(w http.ResponseWriter, r *http.Request){
	var creds Credentials

	//decode user data and save in Credentials struct
	err :=  json.NewDecoder(r.Body).Decode(&creds)
	if err != nil{
		app.errorJSON(w, errors.New("Unauthorized"))
		return
	}

	//we are gonna hash this
	hashedPassword := validUser.Password

	//CompareHashAndPassword compares a bcrypt hashed password with its possible plaintext equivalent.
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword),[]byte(creds.Password))
	if err != nil{
		app.errorJSON(w,errors.New("Unauthorised"))
		return
	}

	//Setting properties to our jwt auth token
	var claims jwt.Claims
	claims.Subject = fmt.Sprint(validUser.ID)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	claims.Issuer = "mydomain.com"
	claims.Audiences = []string{"mydomain.com"}

	//HMACSIGN creats signs in our token and we are using jwt.HS256 algorithm for it and then send our secret []byte in it
	jwtBytes,err := claims.HMACSign(jwt.HS256,[]byte(app.config.jwt.secret))
	if err != nil{
		app.errorJSON(w,errors.New("Unauthorised"))
		return
	}

	//finally sending sending respond to frontend
	app.writeJSON(w,http.StatusOK,jwtBytes,"response")
}