package routes

import (
	"net/http"
	"path/filepath"

	"github.com/ladecadence/MapaLabs/pkg/config"
	"github.com/ladecadence/MapaLabs/pkg/controllers"
	"github.com/ladecadence/MapaLabs/pkg/database"
)

func RegisterRoutes(db database.SQLite, config config.Config, router *http.ServeMux) {
	// web
	router.Handle("GET /static/", http.StripPrefix("/static", http.FileServer(http.Dir(filepath.Join(config.MainPath, "static")))))
	router.HandleFunc("GET /", controllers.ConfMiddleWare(db, config, controllers.WebRoot))
	router.HandleFunc("POST /login", controllers.ConfMiddleWare(db, config, controllers.WebLogin))
	router.HandleFunc("POST /logout", controllers.ConfMiddleWare(db, config, controllers.WebLogout))
	router.HandleFunc("POST /newlab", controllers.ConfMiddleWare(db, config, controllers.WebNewLab))

	// labs
	router.HandleFunc("GET /api/labs", controllers.ConfMiddleWare(db, config, controllers.ApiGetLabs))
	router.HandleFunc("GET /api/lab/{id}", controllers.ConfMiddleWare(db, config, controllers.ApiGetLab))
	router.HandleFunc("POST /api/newlab", controllers.ConfMiddleWare(db, config, controllers.ApiNewLab))
}
