package controllers

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"html"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ladecadence/MapaLabs/pkg/models"
)

func WebLogin(writer http.ResponseWriter, request *http.Request) {
	username := request.FormValue("username")
	password := request.FormValue("password")

	user, err := db.GetUser(username)
	if err != nil {
		log.Printf("❌ User not in DB: %v", err.Error())
		writer.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}
	passwordHash := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
	passwordMatch := (subtle.ConstantTimeCompare([]byte(passwordHash), []byte(user.Password)) == 1)
	if passwordMatch {
		// generate tokens
		sessiontoken, csrftoken := GenTokens()

		http.SetCookie(writer, &http.Cookie{
			Name:     "username",
			Value:    username,
			Expires:  time.Now().Add(60 * time.Minute),
			HttpOnly: true,
		})

		// session token
		http.SetCookie(writer, &http.Cookie{
			Name:     "session_token",
			Value:    sessiontoken,
			Expires:  time.Now().Add(60 * time.Minute),
			HttpOnly: true,
		})

		// and CSRF token
		http.SetCookie(writer, &http.Cookie{
			Name:     "csrf_token",
			Value:    csrftoken,
			Expires:  time.Now().Add(60 * time.Minute),
			HttpOnly: false,
		})

		// and update user
		user.Token = sessiontoken
		user.CSRF = csrftoken
		db.UpsertUser(user)

		// go to main page
		http.Redirect(writer, request, "/", http.StatusSeeOther)

	} else {
		writer.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
	}
}

func WebLogout(writer http.ResponseWriter, request *http.Request) {
	if err := Authorize(request, db); err != nil {
		writer.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// clear cookies
	// user
	http.SetCookie(writer, &http.Cookie{
		Name:     "username",
		Value:    "",
		Expires:  time.Now().Add(-time.Minute),
		HttpOnly: false,
	})

	// token
	http.SetCookie(writer, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Minute),
		HttpOnly: true,
	})

	// and CSRF token
	http.SetCookie(writer, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Minute),
		HttpOnly: false,
	})

	// clear tokens from database
	usercookie, err := request.Cookie("username")
	if err != nil || usercookie.Value == "" {
		writer.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}
	user, err := db.GetUser(usercookie.Value)
	if err != nil {
		writer.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user.Token = ""
	user.CSRF = ""
	db.UpsertUser(user)

	// go to main page
	http.Redirect(writer, request, "/", http.StatusSeeOther)
}

func WebRoot(writer http.ResponseWriter, request *http.Request) {
	tmpl, err := template.ParseFiles(filepath.Join(conf.MainPath, "html/index.html"))

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("Problem loading web page"))
		return
	}

	// labs, err := db.GetLabs()
	// if err != nil {
	// 	writer.WriteHeader(http.StatusInternalServerError)
	// 	writer.Write([]byte("Problem getting labs from database"))
	// 	return
	// }

	// pass user or not
	if err = Authorize(request, db); err != nil {
		log.Printf("❌ Not logged in: %v", err.Error())
		err = tmpl.Execute(writer, nil)
		return
	} else {
		usercookie, err := request.Cookie("username")
		if err != nil || usercookie.Value == "" {
			// what happen!
			err = tmpl.Execute(writer, nil)
			return
		}
		user, err := db.GetUser(usercookie.Value)
		if err != nil {
			// stranger things
			err = tmpl.Execute(writer, nil)
			return
		}

		// ok, render with user
		err = tmpl.Execute(writer, user)

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("Problem rendering web page: "))
			writer.Write([]byte(err.Error()))
			return
		}
	}
}

func WebNewLab(writer http.ResponseWriter, request *http.Request) {
	// check auth
	if err := Authorize(request, db); err != nil {
		writer.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// check CSRF
	if err := CheckCSRF(request, db); err != nil {
		writer.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(writer, "Unauthorized (CSRF)", http.StatusUnauthorized)
		return
	}

	// ok, get data
	file, handler, err := request.FormFile("image")
	if err != nil {
		log.Printf("❌ Error retrieving the file: %v", err.Error())
		return
	}
	defer file.Close()

	f, err := os.OpenFile(conf.ImagePath+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("❌ Error saving the file: %v", err.Error())
		return
	}
	defer f.Close()

	io.Copy(f, file)
	log.Printf("✅ Image file saved")

	// try to create new mission
	lab := models.Lab{}
	lab.Name = html.EscapeString(request.FormValue("name"))
	lab.City = html.EscapeString(request.FormValue("city"))
	lab.Country = html.EscapeString(request.FormValue("country"))
	lab.Description = html.EscapeString(request.FormValue("description"))
	lab.Date = html.EscapeString(request.FormValue("date"))
	lab.Works = html.EscapeString(request.FormValue("works"))
	lab.Motivations = html.EscapeString(request.FormValue("motivations"))
	lab.Networks = html.EscapeString(request.FormValue("networks"))
	lab.Web = html.EscapeString(request.FormValue("web"))
	lab.Mastodon = request.FormValue("mastodon")
	lab.Instagram = html.EscapeString(request.FormValue("instagram"))
	lab.Facebook = html.EscapeString(request.FormValue("facebook"))
	lab.Twitter = html.EscapeString(request.FormValue("twitter"))
	lab.Spotify = html.EscapeString(request.FormValue("spotify"))
	lab.Linkedin = html.EscapeString(request.FormValue("linkedin"))
	lab.TikTok = html.EscapeString(request.FormValue("tiktok"))
	lab.Twitch = html.EscapeString(request.FormValue("twitch"))
	lab.Flickr = html.EscapeString(request.FormValue("flickr"))
	lab.Youtube = html.EscapeString(request.FormValue("youtube"))
	lab.Delegate = html.EscapeString(request.FormValue("delegate"))
	lab.DelegateDescription = html.EscapeString(request.FormValue("delegate_description"))
	lab.DelegatePosition = html.EscapeString(request.FormValue("delegate_position"))
	lab.Image = conf.ImagePath + handler.Filename
	lab.Latitude, err = strconv.ParseFloat(request.FormValue("latitude"), 64)
	lab.Longitude, err = strconv.ParseFloat(request.FormValue("longitude"), 64)

	if err != nil {
		log.Printf("❌ Error decoding form: %v", err.Error())
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("<h1>Data error</h1>\n"))
		return
	}

	err = db.UpsertLab(lab)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("{}\n"))
		return
	}

	// ok, go to main page
	http.Redirect(writer, request, "/", http.StatusSeeOther)

}
