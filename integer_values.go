package raphanus

import "github.com/msoap/raphanus/common"

// GetInt - get integer value by key
func (db *DB) GetInt(key string) (int64, error) {
	if db.withLock {
		db.RLock()
		defer db.RUnlock()
	}

	rawVal, ok := db.data[key]
	if !ok {
		return 0, raphanuscommon.ErrKeyNotExists
	}

	value, ok := rawVal.(int64)
	if !ok {
		return 0, raphanuscommon.ErrKeyTypeMissmatch
	}

	return value, nil
}

// SetInt - create/update integer value by key
func (db *DB) SetInt(key string, value int64, ttl int) error {
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

// UpdateInt - update integer value by exists key
func (db *DB) UpdateInt(key string, value int64) (err error) {
	if db.withLock {
		db.Lock()
		defer db.Unlock()
	}

	_, ok := db.data[key]
	if !ok {
		return raphanuscommon.ErrKeyNotExists
	}

	db.data[key] = value

	return nil
}

// IncrInt - increment integer value on 1
func (db *DB) IncrInt(key string) (err error) {
	if db.withLock {
		db.Lock()
		defer db.Unlock()
	}

	err = db.addInt(key, 1)
	return err
}

// DecrInt - decrement integer value on 1
func (db *DB) DecrInt(key string) (err error) {
	if db.withLock {
		db.Lock()
		defer db.Unlock()
	}

	err = db.addInt(key, -1)
	return err
}

func (db *DB) addInt(key string, value int64) (err error) {
	_, ok := db.data[key]
	if !ok {
		return raphanuscommon.ErrKeyNotExists
	}

	_, ok = db.data[key].(int64)
	if !ok {
		return raphanuscommon.ErrKeyTypeMissmatch
	}

	item := db.data[key].(int64)
	db.data[key] = item + value

	return nil
}
