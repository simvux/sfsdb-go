package cached

import (
	//"bytes"
	"fmt"
	fs "github.com/AlmightyFloppyFish/sfsdb-go/filesystem"
	//"github.com/vmihailenco/msgpack"
	"sort"
	"sync"
)

// Cache holds copies of key/value stored in ram for faster access
// A key/value pair can never be in cache without being saved on disk,
// however can be saved on disk without being saved in cache
type Cache struct{ content sync.Map }

// CacheCounter keeps track of which keys gets used the most,
// to later decide which pairs should be cached
type CacheCounter struct {
	content map[string]uint64
	lock    *sync.RWMutex
}

// NewCache inits the map
func NewCache() *Cache {
	return &Cache{sync.Map{}}
}

// NewCacheCounter inits the map
func NewCacheCounter() CacheCounter {
	return CacheCounter{
		content: make(map[string]uint64),
		lock:    &sync.RWMutex{},
	}
}

// Add or edit key in cache
func (c *Cache) Add(key string, value []byte) {
	c.content.Store(key, value)
}

// Remove from cache
func (c *Cache) Remove(key string) {
	c.content.Delete(key)
}

func (c *Cache) Load(key string, dest interface{}) error {
	encoded, exists := c.content.Load(key)
	if !exists {
		return fmt.Errorf("Does not exist")
	}

	dest = encoded
	return nil
}

// AddTracker adds a usage tracker
func (cc CacheCounter) AddTracker(key string) {
	cc.lock.Lock()
	cc.content[key] = 0
	cc.lock.Unlock()
}

// DelTracker removes a usage tracker
func (cc CacheCounter) DelTracker(key string) {
	cc.lock.Lock()
	delete(cc.content, key)
	cc.lock.Unlock()
}

// IncreaseUse adds one to the use counter for key
func (cc CacheCounter) IncreaseUse(key string) {
	cc.lock.Lock()
	cc.content[key]++
	cc.lock.Unlock()
}

// Len gets length of the content field
func (cc CacheCounter) len() int {
	cc.lock.RLock()
	length := len(cc.content)
	cc.lock.RUnlock()
	return length
}

// Reset prevents integer overflows by pushing back all values, maintaining the percentual
// difference between the counts.
func (cc CacheCounter) Reset() {
	cc.lock.Lock()
	for i, v := range cc.content {
		cc.content[i] = v / 5
	}
	cc.lock.Unlock()
}

type pair struct {
	key   string
	value uint64
}

// Resync loads the db.cacheLimit top keys into db.cache
func (db *Cached) Resync() {
	db.cacheCount.lock.RLock() // TODO: Am i locking this safely?
	pairs := make([]pair, db.cacheCount.len())
	count := 0
	for i, v := range db.cacheCount.content {
		pairs[count] = pair{i, v}
		count++
	}
	db.cacheCount.lock.RUnlock()
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].value > pairs[j].value })

	for i, pair := range pairs {
		if i > count {
			break
		} else if uint64(i) >= db.cacheLimit && db.cacheLimit > 0 {
			break
		}
		if _, exists := db.cache.content.Load(pair.key); exists {
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
		db.cache.content.Store(pair.key, data)
	}
}
