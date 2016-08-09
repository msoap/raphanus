package raphanus

import "github.com/msoap/raphanus/common"

// ListValue - list value type
type ListValue []string

// GetList - get list value by key
func (db *DB) GetList(key string) (value ListValue, err error) {
	if db.withLock {
		db.RLock()
		defer db.RUnlock()
	}

	rawVal, ok := db.data[key]
	if !ok {
		return value, raphanuscommon.ErrKeyNotExists
	}

	value, ok = rawVal.val.(ListValue)
	if !ok {
		return value, raphanuscommon.ErrKeyTypeMissmatch
	}

	return value, err
}

// SetList - create/update list value by key
func (db *DB) SetList(key string, value ListValue, ttl int) {
	if db.withLock {
		db.Lock()
		defer db.Unlock()
	}

	item := db.data[key]
	item.val = value
	db.data[key] = item

	db.setTTL(key, ttl)

	return
}

// UpdateList - update list value by exists key
func (db *DB) UpdateList(key string, value ListValue) (err error) {
	if db.withLock {
		db.Lock()
		defer db.Unlock()
	}

	if _, ok := db.data[key]; !ok {
		return raphanuscommon.ErrKeyNotExists
	}

	item := db.data[key]
	item.val = value
	db.data[key] = item

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

	valueList, ok := rawVal.val.(ListValue)
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

	valueList, ok := rawVal.val.(ListValue)
	if !ok {
		return raphanuscommon.ErrKeyTypeMissmatch
	}

	if index < 0 || index >= len(valueList) {
		return raphanuscommon.ErrListOutOfRange
	}

	db.data[key].val.(ListValue)[index] = value
	return nil
}
