package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/ladecadence/MapaLabs/pkg/database"
)

var AuthError = errors.New("Unauthorized")

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
