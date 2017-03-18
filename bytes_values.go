package raphanus

import "github.com/msoap/raphanus/common"

// GetBytes - get []byte value by key
func (db *DB) GetBytes(key string) ([]byte, error) {
	if db.withLock {
		db.RLock()
		defer db.RUnlock()
	}

	rawVal, ok := db.data[key]
	if !ok {
		return nil, raphanuscommon.ErrKeyNotExists
	}

	value, ok := rawVal.([]byte)
	if !ok {
		return nil, raphanuscommon.ErrKeyTypeMissmatch
	}

	return value, nil
}

// SetBytes - create/update []byte value by key
func (db *DB) SetBytes(key string, value []byte, ttl int) error {
	if db.withLock {
		db.Lock()
		defer db.Unlock()
	}

	if !isValidKey(key) {
		return raphanuscommon.ErrKeyIsNotValid
	}

	db.data[key] = value
	db.setTTL(key, ttl)

	return nil
}

// UpdateBytes - update []byte value by exists key
func (db *DB) UpdateBytes(key string, value []byte) (err error) {
	if db.withLock {
		db.Lock()
		defer db.Unlock()
	}

	if _, ok := db.data[key]; !ok {
		return raphanuscommon.ErrKeyNotExists
	}

	db.data[key] = value

	return nil
}
