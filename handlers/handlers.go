// Package handlers provides HTTP handlers of under_pressure service
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
			err := storage.AddPressure(input)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	}
}

// ListPressures lists all pressure points
func ListPressures(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		pressures, err := storage.ListPressures()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			encoder := json.NewEncoder(w)
			err := encoder.Encode(pressures)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	}
}

// ListPressuresCSV lists all pressure points in CSV format
func ListPressuresCSV(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		pressures, err := storage.ListPressures()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", "text/csv")
			fmt.Fprintf(w, "Date,Systolic,Diastolic,Heart Rate\n")
			for _, p := range pressures {
				fmt.Fprintf(w, "%s,%d,%d,%d\n", time.Unix(p.Timestamp, 0).Format("2006/01/02 15:04:05"), p.Systolic, p.Diastolic, p.HeartRate)
			}
		}
	}
}

// RemovePressure removes a pressure point key (timestamp)
func RemovePressure(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		timestamp := r.URL.Query().Get("timestamp")
		if len(timestamp) != 0 {
			err := storage.DeletePressure(timestamp)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}
