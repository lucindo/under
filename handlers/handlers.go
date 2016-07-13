// Package handlers provides HTTP handlers of under_pressure service
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/lucindo/under_pressure/log"
	"github.com/lucindo/under_pressure/pressure"
)

// PostPressure function inserts a new Pressure point on the storage
func PostPressure(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		var input pressure.Pressure

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&input)

		if err != nil || !input.Valid() {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			log.Logger.Println(input)
		}
	}
}
