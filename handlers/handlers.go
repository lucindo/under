// Package handlers provides HTTP handlers of under_pressure service
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/lucindo/under_pressure/pressure"
	"github.com/lucindo/under_pressure/storage"
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
			storage.AddPressure(input)
		}
	}
}

// ListPressures lists all pressure points
func ListPressures(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		pressures := storage.ListPressures()
		encoder := json.NewEncoder(w)
		err := encoder.Encode(pressures)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
