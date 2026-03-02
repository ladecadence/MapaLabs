package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ladecadence/MapaLabs/pkg/color"
	"github.com/ladecadence/MapaLabs/pkg/config"
	"github.com/ladecadence/MapaLabs/pkg/database"
	"github.com/ladecadence/MapaLabs/pkg/routes"
)

func main() {
	// command line flags
	configFile := flag.String("c", "config.toml", "Config file")
	flag.Parse()

	// try to load config file, if not, create a defualt configuration file
	// at standard config directory
	conf, err := config.GetConfig(*configFile)
	if err != nil {
		fmt.Printf("Can't open configuration file %s\n", *configFile)
		os.Exit(1)
	}

	// open and init database
	database := &database.SQLite{}
	_, err = database.Open(conf.Database)
	if err != nil {
		panic("Error opening DB: " + err.Error())
	}
	err = database.Init()
	if err != nil {
		panic("Error initializing DB: " + err.Error())
	}

	// create server mux and routes
	mux := http.NewServeMux()
	routes.RegisterRoutes(*database, conf, mux)

	const logo = `
	
	███╗   ███╗ █████╗ ██████╗  █████╗ ██╗      █████╗ ██████╗ ███████╗
	████╗ ████║██╔══██╗██╔══██╗██╔══██╗██║     ██╔══██╗██╔══██╗██╔════╝
	██╔████╔██║███████║██████╔╝███████║██║     ███████║██████╔╝███████╗
	██║╚██╔╝██║██╔══██║██╔═══╝ ██╔══██║██║     ██╔══██║██╔══██╗╚════██║
	██║ ╚═╝ ██║██║  ██║██║     ██║  ██║███████╗██║  ██║██████╔╝███████║
	╚═╝     ╚═╝╚═╝  ╚═╝╚═╝     ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝╚═════╝ ╚══════╝
                                                                   
	`

	fmt.Fprintf(os.Stderr, "%s", color.Purple+logo+color.Reset)
	log.Printf("🌍 "+color.Green+"MapaLabs version "+color.Purple+"%s"+color.Green+" listening on port "+color.Yellow+"%d"+color.Reset, conf.Version, conf.Port)

	// launch server
	addr := fmt.Sprintf("%s:%d", conf.Addr, conf.Port)
	log.Fatal(http.ListenAndServe(addr, mux))

}
