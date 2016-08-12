package raphanus

import "github.com/msoap/raphanus/common"

// GetList - get list value by key
func (db *DB) GetList(key string) (value raphanuscommon.ListValue, err error) {
	if db.withLock {
		db.RLock()
		defer db.RUnlock()
	}

	rawVal, ok := db.data[key]
	if !ok {
		return value, raphanuscommon.ErrKeyNotExists
	}

	value, ok = rawVal.(raphanuscommon.ListValue)
	if !ok {
		return value, raphanuscommon.ErrKeyTypeMissmatch
	}

	return value, err
}

// SetList - create/update list value by key
func (db *DB) SetList(key string, value raphanuscommon.ListValue, ttl int) error {
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

// UpdateList - update list value by exists key
func (db *DB) UpdateList(key string, value raphanuscommon.ListValue) (err error) {
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

// GetListItem - get one item from list value by key
func (db *DB) GetListItem(key string, index int) (string, error) {
	if db.withLock {
		db.RLock()
		defer db.RUnlock()
	}

	rawVal, ok := db.data[key]
	if !ok {
		return "", raphanuscommon.ErrKeyNotExists
	}

	valueList, ok := rawVal.(raphanuscommon.ListValue)
	if !ok {
		return "", raphanuscommon.ErrKeyTypeMissmatch
	}

	if index < 0 || index >= len(valueList) {
		return "", raphanuscommon.ErrListOutOfRange
	}

	result := valueList[index]
	return result, nil
}

// SetListItem - set one item of list value by key
func (db *DB) SetListItem(key string, index int, value string) error {
	if db.withLock {
		db.RLock()
		defer db.RUnlock()
	}

	rawVal, ok := db.data[key]
	if !ok {
		return raphanuscommon.ErrKeyNotExists
	}

	valueList, ok := rawVal.(raphanuscommon.ListValue)
	if !ok {
		return raphanuscommon.ErrKeyTypeMissmatch
	}

	if index < 0 || index >= len(valueList) {
		return raphanuscommon.ErrListOutOfRange
	}

	db.data[key].(raphanuscommon.ListValue)[index] = value
	return nil
}
