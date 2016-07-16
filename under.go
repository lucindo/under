package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/lucindo/under_pressure/handlers"
	"github.com/lucindo/under_pressure/log"
	"github.com/lucindo/under_pressure/storage"
)

func main() {
	port := flag.Int("port", 8080, "port number")
	dbfile := flag.String("dbfile", "/tmp/pressure.bolt", "database filename")
	staticDir := flag.String("staticdir", "./static", "static files dir path")
	help := flag.Bool("help", false, "prints usage and exits")
	interrupt := make(chan os.Signal, 1)

	log.Init()
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	signal.Notify(interrupt, os.Interrupt)
	go func() {
		for sig := range interrupt {
			log.Logger.Printf("exiting due signal: %s", sig)
			// TODO: stop the server gracefully
			os.Exit(0)
		}
	}()

	storage.Init(*dbfile)
	defer storage.Close()

	http.Handle("/", http.FileServer(http.Dir(*staticDir)))
	http.HandleFunc("/new", handlers.PostPressure)
	http.HandleFunc("/all", handlers.ListPressures)
	http.HandleFunc("/all.csv", handlers.ListPressuresCSV)

	server := &http.Server{
		Addr:     fmt.Sprintf(":%d", *port),
		ErrorLog: log.Logger,
	}
	log.Logger.Printf("listening on port %d (static files from %s)\n", *port, *staticDir)
	log.Logger.Fatal(server.ListenAndServe())
}
