package cached

import (
	"bytes"
	"fmt"
	fs "github.com/AlmightyFloppyFish/sfsdb-go/filesystem"
	"github.com/vmihailenco/msgpack"
	"sort"
)

// Cache holds copies of key/value stored in ram for faster access
// A key/value pair can never be in cache without being saved on disk,
// however can be saved on disk without being saved in cache
type Cache map[string][]byte

// CacheCounter keeps track of which keys gets used the most,
// to later decide which pairs should be cached
type CacheCounter map[string]uint64

// NewCache inits the map
func NewCache() Cache {
	return Cache(make(map[string][]byte))
}

// NewCacheCounter inits the map
func NewCacheCounter() CacheCounter {
	return CacheCounter(make(map[string]uint64))
}

// Add or edit key in cache
func (c Cache) Add(key string, value []byte) {
	c[key] = value
}

// Remove from cache
func (c Cache) Remove(key string) {
	delete(c, key)
}

func (c Cache) Load(key string, dest interface{}) bool {
	encoded, exists := c[key]
	if !exists {
		return false
	}
	reader := bytes.NewReader(encoded)
	dec := msgpack.NewDecoder(reader)

	if err := dec.Decode(dest); err != nil {
		fmt.Printf("\nsfsdb: Cache violation (%s): %s", key, err.Error())
		return false
	}
	return true
}

// AddTracker adds a usage tracker
func (cc CacheCounter) AddTracker(key string) {
	cc[key] = 0
}

// DelTracker removes a usage tracker
func (cc CacheCounter) DelTracker(key string) {
	delete(cc, key)
}

// IncreaseUse adds one to the use counter for key
func (cc CacheCounter) IncreaseUse(key string) {
	cc[key]++
}

// Reset prevents integer overflows by pushing back all values, maintaining the percentual
// difference between the counts.
func (cc CacheCounter) Reset() {
	for i, v := range cc {
		cc[i] = v / 5
	}
}

type pair struct {
	key   string
	value uint64
}

// Resync loads the db.cacheLimit top keys into db.cache
func (db *Cached) Resync() {
	pairs := make([]pair, len(db.cacheCount))
	count := 0
	for i, v := range db.cacheCount {
		pairs[count] = pair{i, v}
		count++
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].value > pairs[j].value })

	for i, pair := range pairs {
		if i > count {
			break
		} else if uint64(i) >= db.cacheLimit && db.cacheLimit > 0 {
			break
		}
		if _, exists := db.cache[pair.key]; exists {
			// Skip this one because it already exists
			continue
		}

		path := fs.NewFilepath(db.Location())
		path.Append(pair.key)
		data, err := fs.LoadRaw(path)
		if err != nil {
			fmt.Printf("\nsfsdb: File and Cache mismatch (%s): %s", pair.key, err.Error())
			continue
		}
		db.cache[pair.key] = data
	}
}
