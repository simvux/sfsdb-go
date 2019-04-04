package sfsdb

import (
	"github.com/AlmightyFloppyFish/sfsdb-go/cached"
	"github.com/AlmightyFloppyFish/sfsdb-go/simple"
	"os"
)

// Database .
type Database interface {
	// Get the root directory of the database
	Location() string

	// Check if specified key exists
	Exists(key string) bool

	// Save value to key
	Save(key string, value interface{}) error

	// Load value stored in key to destination (must be pointer)
	Load(key string, dest interface{}) error

	// Delete the key/value pair
	Delete(string) error
}

// New database, Set cache to 0 for uncached database.
func New(location string, cache uint64) Database {
	os.MkdirAll(location, 0777)

	if cache > 0 {
		db := cached.New(location, cache)
		return &db
	}
	db := simple.New(location)
	return &db
}
