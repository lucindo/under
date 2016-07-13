package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/lucindo/under_pressure/handlers"
	"github.com/lucindo/under_pressure/log"
)

func main() {
	port := flag.Int("port", 8080, "port number")

	log.Init()
	flag.Parse()

	http.HandleFunc("/new", handlers.PostPressure)

	log.Logger.Printf("listening on port %d\n", *port)
	log.Logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
