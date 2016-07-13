// Package log provides logging facilities
package log

import (
	"log"
	"os"
)

// Logger object
var Logger *log.Logger

// Init function constructs the global log object
func Init() {
	Logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
}
