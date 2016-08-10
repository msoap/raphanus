package raphanus

import "github.com/msoap/raphanus/common"

// GetStr - get string value by key
func (db *DB) GetStr(key string) (string, error) {
	if db.withLock {
		db.RLock()
		defer db.RUnlock()
	}

	rawVal, ok := db.data[key]
	if !ok {
		return "", raphanuscommon.ErrKeyNotExists
	}

	value, ok := rawVal.(string)
	if !ok {
		return "", raphanuscommon.ErrKeyTypeMissmatch
	}

	return value, nil
}

// SetStr - create/update string value by key
func (db *DB) SetStr(key, value string, ttl int) {
	if db.withLock {
		db.Lock()
		defer db.Unlock()
	}

	db.data[key] = value
	db.setTTL(key, ttl)

	return
}

// UpdateStr - update string value by exists key
func (db *DB) UpdateStr(key, value string) (err error) {
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
