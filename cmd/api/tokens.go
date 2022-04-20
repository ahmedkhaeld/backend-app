package main

import (
	"backend/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pascaldekloe/jwt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// validUser a dummy user
var validUser = models.User{
	ID:       10,
	Email:    "me@here.com",
	Password: "$2a$12$YZmO3zxVXaKGXORRDxMleOD8COPtz85eSfuxB3ulSwfZmQ6uNzmE2"}

// Credentials to test whether someone is allowed to log into the system or not
// against their email and password
type Credentials struct {
	Username string `json:"email"`
	Password string `json:"password"`
}

func (app *application) SignIn(w http.ResponseWriter, r *http.Request) {

	var creds Credentials

	// read the json body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"), 400)
		return
	}

	// this step is a shortcut for get username&password&email by query the db
	// to check if user matches that username, password hash stored in the db against
	// what is being provided
	hashedPassword := validUser.Password

	// check the hash for the password matches the hash in the db
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password))
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"))
		return
	}
	// until this point we have a valid user and password

	// create our jwt and send it back
	// make the claims
	var claims jwt.Claims
	claims.Subject = fmt.Sprint(validUser.ID)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	claims.Issuer = "mydomain.com"
	claims.Audiences = []string{"mydomain.com"}

	//create the tokens
	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secret))
	if err != nil {
		app.errorJSON(w, errors.New("error signing"))
		return
	}

	app.writeJSON(w, http.StatusOK, string(jwtBytes), "response")

}
