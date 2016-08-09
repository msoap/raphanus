/*
Package raphanus - simple implementation of Redis-like in-memory cache
*/
package raphanus

import (
	"errors"
	"sync"
	"time"
)

// MaxStringValueLength - string value max length
const MaxStringValueLength = 1024 * 1024 // 1 MB

// Common errors
var (
	ErrKeyNotExists     = errors.New("Key not exists")
	ErrKeyTypeMissmatch = errors.New("The type does not match")
	ErrListOutOfRange   = errors.New("List index is out of range")
	ErrDictKeyNotExists = errors.New("Dict, key not exists")
	ErrDictKeyIsEmpty   = errors.New("Key or dict key is empty")
	ErrTTLIsntCorrect   = errors.New("TTL parameter isn't correct")
)

type value struct {
	val interface{}
	ttl time.Time
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

// setTTL - set TTL on one value in DB
func (db *DB) setTTL(key string, ttl int) {
	if ttl > 0 {
		go func() {
			time.Sleep(time.Duration(ttl) * time.Second)
			db.Lock()
			defer db.Unlock()

			delete(db.data, key)
		}()
	}
}
