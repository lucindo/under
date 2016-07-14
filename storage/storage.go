// Package storage provides helper functions to read/write pressure data points
package storage

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/lucindo/under_pressure/log"
	"github.com/lucindo/under_pressure/pressure"
)

// private boltdb vars
var (
	db     *bolt.DB
	bucket = "pressure"
)

// Init initializes the database
func Init(dbFile string) {
	var err error

	log.Logger.Printf("openning database file: %s\n", dbFile)
	db, err = bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Logger.Printf("error opening database file [%s]: %v\n", dbFile, err)
		panic(err)
	}

	log.Logger.Printf("creating bucket [%s] (if not exists)\n", dbFile)
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Logger.Printf("error creating bucket [%s]: %v\n", bucket, err)
		panic(err)
	}
}

// Close the databse
func Close() {
	log.Logger.Printf("closing database\n")
	err := db.Close()
	if err != nil {
		log.Logger.Printf("error closing database: %v", err)
	}
}

// AddPressure adds a new pressure datapoint to the database
func AddPressure(p pressure.Pressure) {
	log.Logger.Printf("adding new pressure point to the database: %s\n", p)
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		encoded, err := json.Marshal(p)
		if err != nil {
			return err
		}
		return b.Put([]byte(string(p.Timestamp)), encoded)
	})
	if err != nil {
		log.Logger.Printf("error adding %s to database: %v\n", p, err)
	}
}

// ListPressures lists all pressure points on database
func ListPressures() []pressure.Pressure {
	var pList []pressure.Pressure
	log.Logger.Printf("getting all pressure points")
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.ForEach(func(k, v []byte) error {
			var decoded pressure.Pressure
			err := json.Unmarshal(v, &decoded)
			if err != nil {
				log.Logger.Printf("error decoding pressure point key [%v]: %v\n", k, err)
			}
			pList = append(pList, decoded)
			return nil
		})
		return err
	})
	if err != nil {
		log.Logger.Printf("error reading database: %v\n", err)
	}
	return pList
}
