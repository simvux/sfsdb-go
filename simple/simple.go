package simple

import (
	fs "github.com/AlmightyFloppyFish/sfsdb-go/filesystem"
	"os"
)

type Simple struct {
	location string
}

func New(location string) Simple {
	return Simple{
		location: location,
	}
}

func (db *Simple) Location() string {
	return db.location
}

func (db *Simple) Save(key string, data interface{}) error {
	path := fs.NewFilepath(db.Location())
	path.Append(key)
	return fs.Save(path, data)
}

func (db *Simple) Load(key string, dest interface{}) error {
	path := fs.NewFilepath(db.Location())
	path.Append(key)
	return fs.Load(path, dest)
}

func (db *Simple) Exists(key string) bool {
	path := fs.NewFilepath(db.Location())
	path.Append(key)
	_, err := os.Stat(path.Unwrap())
	return !os.IsNotExist(err)
}

func (db *Simple) Delete(key string) error {
	path := fs.NewFilepath(db.Location())
	path.Append(key)
	return fs.Delete(path)
}
