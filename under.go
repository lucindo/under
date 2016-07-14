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
	interrupt := make(chan os.Signal, 1)

	log.Init()
	flag.Parse()

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

	http.HandleFunc("/new", handlers.PostPressure)
	http.HandleFunc("/all", handlers.ListPressures)

	server := &http.Server{
		Addr:     fmt.Sprintf(":%d", *port),
		ErrorLog: log.Logger,
	}
	log.Logger.Printf("listening on port %d\n", *port)
	log.Logger.Fatal(server.ListenAndServe())
}
