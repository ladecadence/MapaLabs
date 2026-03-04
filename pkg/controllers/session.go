package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/ladecadence/MapaLabs/pkg/database"
)

var AuthError = errors.New("Unauthorized")
var CSRFError = errors.New("CSRF Token invalid")

func Authorize(request *http.Request, db database.SQLite) error {
	// get the user from the cookie
	usercookie, err := request.Cookie("username")
	if err != nil || usercookie.Value == "" {
		log.Printf("❌ No user cookie %v", err.Error())
		return AuthError
	}
	user, err := db.GetUser(usercookie.Value)
	if err != nil {
		log.Printf("❌ User not in DB %v", err.Error())
		return AuthError
	}

	// get the session token from the cookie
	st, err := request.Cookie("session_token")
	if err != nil || st.Value == "" || st.Value != user.Token {
		log.Printf("❌ No session token")
		return AuthError
	}

	return nil
}

func CheckCSRF(request *http.Request, db database.SQLite) error {
	// get the user from the cookie
	usercookie, err := request.Cookie("username")
	if err != nil || usercookie.Value == "" {
		log.Printf("❌ No user cookie %v", err.Error())
		return CSRFError
	}
	user, err := db.GetUser(usercookie.Value)
	if err != nil {
		log.Printf("❌ User not in DB %v", err.Error())
		return CSRFError
	}

	// get the CSRF token from the form and test it
	if csrf := request.FormValue("csrf"); csrf == "" || csrf != user.CSRF {
		log.Printf("❌ CSRF Token not present or invalid")
		return CSRFError
	}

	return nil
}
