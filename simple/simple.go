package simple

import (
	"fmt"
	fs "github.com/AlmightyFloppyFish/sfsdb-go/filesystem"
	"github.com/AlmightyFloppyFish/sfsdb-go/lock"
	"os"
)

const (
	errIllegalPath   = "Error: Illegal Path %s"
	cacheResyncEvery = 100
)

type Simple struct {
	location string
	locker   *lock.WriteLock
}

func New(location string) Simple {
	return Simple{
		location: location,
		locker:   lock.New(),
	}
}

func (db *Simple) Location() string {
	return db.location
}

func (db *Simple) Save(key string, data interface{}) error {
	path := fs.NewFilepath(db.Location())
	path.Append(key)
	if path.Unwrap() == db.Location() {
		return fmt.Errorf(errIllegalPath, path.Unwrap())
	}

	lock := db.locker.Get(key)
	err := fs.Save(path, data)
	db.locker.Done(key, lock)
	return err
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

	lock := db.locker.Get(key)
	err := fs.Delete(path)
	db.locker.Done(key, lock)
	return err
}
