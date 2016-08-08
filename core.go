/*
Package raphanus - simple implementation of Redis-like in-memory cache
*/
package raphanus

import (
	"errors"
	"sync"
)

var (
	ErrKeyNotExists     = errors.New("Key not exists")
	ErrKeyTypeMissmatch = errors.New("The type does not match")
)

type value struct {
	val interface{}
	ttl int
}

type DB struct {
	data     map[string]value
	withLock bool // execute get/set methods under lock
	*sync.RWMutex
}

func New() DB {
	return DB{
		data:     map[string]value{},
		withLock: true,
		RWMutex:  new(sync.RWMutex),
	}
}

func (db *DB) UnderRLock(fn func()) {
	db.RLock()
	db.withLock = false
	fn()
	db.RUnlock()
	db.withLock = true
}

func (db *DB) UnderLock(fn func()) {
	db.Lock()
	db.withLock = false
	fn()
	db.Unlock()
	db.withLock = true
}

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

func (db *DB) Len() int {
	return len(db.data)
}
