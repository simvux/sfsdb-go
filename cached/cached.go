package cached

import (
	"fmt"
	fs "github.com/AlmightyFloppyFish/sfsdb-go/filesystem"
	"os"
)

const (
	errIllegalPath   = "Error: Illegal Path %s"
	cacheResyncEvery = 100
)

// Cached database unlike Simple database holds
// cacheLimit amount of copies of key/value pairs in memory
type Cached struct {
	location   string
	cacheLimit uint64
	cacheTimer uint8
	cacheCount CacheCounter
	cache      Cache
}

// New initializes the database
func New(location string, cache uint64) Cached {
	return Cached{
		location:   location,
		cacheLimit: cache,
		cacheCount: NewCacheCounter(),
		cache:      NewCache(),
	}
}

func (db *Cached) Location() string {
	return db.location
}

func (db *Cached) Save(key string, data interface{}) error {
	db.cacheCount.AddTracker(key)
	path := fs.NewFilepath(db.Location())
	path.Append(key)
	if path.Unwrap() == db.Location() {
		return fmt.Errorf(errIllegalPath, path.Unwrap())
	}
	return fs.Save(path, data)
}

func (db *Cached) Load(key string, dest interface{}) error {
	db.cacheCount.IncreaseUse(key)
	if db.cacheTimer > cacheResyncEvery {
		db.Resync()
		db.cacheTimer = 0
	}
	db.cacheTimer++
	if existed := db.cache.Load(key, dest); existed {
		return nil
	}

	path := fs.NewFilepath(db.Location())
	path.Append(key)
	if path.Unwrap() == db.Location() {
		return fmt.Errorf(errIllegalPath, path.Unwrap())
	}
	return fs.Load(path, dest)
}

func (db *Cached) Exists(key string) bool {
	path := fs.NewFilepath(db.Location())
	path.Append(key)
	_, err := os.Stat(path.Unwrap())
	return os.IsExist(err)
}

func (db *Cached) Delete(key string) error {
	db.cache.Remove(key)
	path := fs.NewFilepath(db.Location())
	path.Append(key)
	return fs.Delete(path)
}
