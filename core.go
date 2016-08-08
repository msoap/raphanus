/*
Package raphanus - simple implementation of Redis-like in-memory cache
*/
package raphanus

import (
	"errors"
	"sync"
)

// Common errors
var (
	ErrKeyNotExists     = errors.New("Key not exists")
	ErrKeyTypeMissmatch = errors.New("The type does not match")
)

type value struct {
	val interface{}
	ttl int
}

// DB - in-memory cache object
type DB struct {
	data     map[string]value
	withLock bool // execute get/set methods under lock
	*sync.RWMutex
}

// New - get new cache object
func New() DB {
	return DB{
		data:     map[string]value{},
		withLock: true,
		RWMutex:  new(sync.RWMutex),
	}
}

// UnderRLock - execute few read-methods undef one RLock
func (db *DB) UnderRLock(fn func()) {
	db.RLock()
	db.withLock = false
	fn()
	db.RUnlock()
	db.withLock = true
}

// UnderLock - execute few RW-methods undef one Lock
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
		return ErrKeyNotExists
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
