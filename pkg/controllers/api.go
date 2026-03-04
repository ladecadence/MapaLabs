package controllers

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"log"

	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/ladecadence/MapaLabs/pkg/config"
	"github.com/ladecadence/MapaLabs/pkg/database"
	"github.com/ladecadence/MapaLabs/pkg/models"
)

var conf config.Config
var db database.SQLite

func ConfMiddleWare(dtb database.SQLite, c config.Config, h http.HandlerFunc) http.HandlerFunc {
	conf = c
	db = dtb
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

func GenTokens() (string, string) {
	sessiontoken := make([]byte, 32)
	csrftoken := make([]byte, 32)
	// can't fail
	rand.Read(sessiontoken)
	rand.Read(csrftoken)

	return base64.URLEncoding.EncodeToString(sessiontoken), base64.URLEncoding.EncodeToString(csrftoken)
}

func CheckBasicAuth(r *http.Request) bool {
	username, password, ok := r.BasicAuth()
	if ok {
		// get username from DB
		user, err := db.GetUser(username)
		if err != nil {
			return false
		}
		// ok, we have a username, check password
		passwordHash := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
		passwordMatch := (subtle.ConstantTimeCompare([]byte(passwordHash), []byte(user.Password)) == 1)
		if passwordMatch {
			return true
		} else {
			return false
		}

	} else {
		return false
	}
}

func ApiGetLabs(writer http.ResponseWriter, request *http.Request) {
	labs, err := db.GetLabs()

	if err != nil || labs == nil {
		writer.WriteHeader(http.StatusNoContent)
		writer.Write([]byte(`{}\n`))
		return
	}

	res, _ := json.Marshal(labs)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(res)
	writer.Write([]byte("\n"))
}

func ApiGetLab(writer http.ResponseWriter, request *http.Request) {
	// get ID
	sid := request.PathValue("id")
	if sid == "" {
		writer.WriteHeader(http.StatusNoContent)
		writer.Write([]byte(`{}\n`))
		return
	}
	id, err := strconv.Atoi(sid)
	if err != nil {
		writer.WriteHeader(http.StatusNoContent)
		writer.Write([]byte(`{}\n`))
		return

	}
	lab, err := db.GetLab(id)

	if err != nil {
		writer.WriteHeader(http.StatusNoContent)
		writer.Write([]byte(`{}\n`))
		return
	}

	res, _ := json.Marshal(lab)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(res)
	writer.Write([]byte("\n"))
}

func ApiNewLab(writer http.ResponseWriter, request *http.Request) {
	// check auth
	authOk := CheckBasicAuth(request)

	if authOk {
		file, handler, err := request.FormFile("image")
		if err != nil {
			fmt.Fprintf(writer, "Error retrieving the file: %v", err)
			return
		}
		defer file.Close()

		f, err := os.OpenFile(conf.ImagePath+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Fprintf(writer, "Error saving the file: %v", err)
			return
		}
		defer f.Close()

		io.Copy(f, file)
		fmt.Fprintf(writer, "File uploaded successfully: %v", handler.Filename)

		// try to create new mission
		lab := models.Lab{}
		lab.Name = request.FormValue("name")
		lab.City = request.FormValue("city")
		lab.Country = request.FormValue("country")
		lab.Description = request.FormValue("description")
		lab.Date = request.FormValue("date")
		lab.Works = request.FormValue("works")
		lab.Motivations = request.FormValue("motivations")
		lab.Networks = request.FormValue("networks")
		lab.Web = request.FormValue("web")
		lab.Mastodon = request.FormValue("mastodon")
		lab.Instagram = request.FormValue("instagram")
		lab.Facebook = request.FormValue("facebook")
		lab.Twitter = request.FormValue("twitter")
		lab.Spotify = request.FormValue("spotify")
		lab.Linkedin = request.FormValue("linkedin")
		lab.TikTok = request.FormValue("tiktok")
		lab.Twitch = request.FormValue("twitch")
		lab.Flickr = request.FormValue("flickr")
		lab.Youtube = request.FormValue("youtube")
		lab.Delegate = request.FormValue("delegate")
		lab.DelegateDescription = request.FormValue("delegate_description")
		lab.DelegatePosition = request.FormValue("delegate_position")
		lab.Image = conf.ImagePath + handler.Filename
		lab.Latitude, err = strconv.ParseFloat(request.FormValue("latitude"), 64)
		lab.Longitude, err = strconv.ParseFloat(request.FormValue("longitude"), 64)

		if err != nil {
			log.Printf("❌ Error decoding form: %v", err.Error())
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("{}\n"))
			return
		}

		err = db.UpsertLab(lab)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("{}\n"))
			return
		}

		data, _ := json.Marshal(lab)
		writer.WriteHeader(http.StatusOK)
		writer.Write(data)
		writer.Write([]byte("\n"))
	} else {
		writer.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
	}
}
