/*
Package raphanus - simple implementation of in-memory cache
*/
package raphanus

import (
	"log"
	"sync"
	"time"

	"github.com/msoap/raphanus/common"
)

// DB - in-memory cache object
type DB struct {
	data              map[string]interface{}
	withLock          bool   // execute get/set methods under lock
	fsStorageName     string // file name for save/load
	fsStorageSyncTime int    // seconds between saves on disk
	*sync.RWMutex
}

// New - get new cache object
func New() DB {
	db := DB{
		data:     map[string]interface{}{},
		withLock: true,
		RWMutex:  new(sync.RWMutex),
	}

	return db
}

// SetStorage - setup persistent storage
func (db *DB) SetStorage(fsStorageName string, fsStorageSyncTime int) DB {
	if fsStorageName == "" {
		log.Fatal("storage file name cannot be empty")
	}

	err := db.fsLoad()
	if err != nil {
		log.Printf("Load from file failed: %s", err)
	}

	if fsStorageSyncTime > 0 {
		go db.fsHandle()
	}

	return *db
}

// UnderLock - execute few RW-methods under one Lock
func (db *DB) UnderLock(fn func()) {
	db.Lock()
	db.withLock = false
	fn()
	db.Unlock()
	db.withLock = true
}

// Remove - remove key from cache
func (db *DB) Remove(key string) (err error) {
	if db.withLock {
		db.Lock()
		defer db.Unlock()
	}

	_, ok := db.data[key]
	if !ok {
		return raphanuscommon.ErrKeyNotExists
	}

	delete(db.data, key)
	return nil
}

// Keys - get all keys from cache (as []string)
func (db *DB) Keys() []string {
	if db.withLock {
		db.RLock()
		defer db.RUnlock()
	}

	result := make([]string, 0, len(db.data))
	for key := range db.data {
		result = append(result, key)
	}

	return result
}

// Len - return length of cache
func (db *DB) Len() int {
	return len(db.data)
}

// setTTL - set TTL on one value in DB
func (db *DB) setTTL(key string, ttl int) {
	if ttl > 0 {
		// race condition if key deleted and created again before fire the ttl event
		time.AfterFunc(time.Duration(ttl)*time.Second, func() {
			db.Lock()
			delete(db.data, key)
			db.Unlock()
		})
	}
}

func isValidKey(key string) bool {
	return len(key) > 0
}
