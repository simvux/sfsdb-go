package main

import (
	"fmt"
	fs "github.com/AlmightyFloppyFish/sfsdb-go/filesystem"
	"os"
)

type Indexed struct {
	location string
	index    Index
}

func (db *Indexed) New(location string, index uint64) {
	// TODO: Init the cache from db.index.location to db.index.loaded
	return Indexed{
		location: location,
		index: Index{
			loaded:   make(map[string]interface{}),
			location: location + "/__INDEX__",
		},
	}
}

func (db *Indexed) Location() {
	return db.location
}

func (db *Indexed) Save() {

}
func (db *Indexed) Load() {

}

func (db *Indexed) Exists() {

}

func (db *Indexed) Delete() {
	// Remember to delete index!
}
