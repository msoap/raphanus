package raphanus

// DictValue - dict value type
type DictValue map[string]string

// GetDict - get dict value by key
func (db *DB) GetDict(key string) (value DictValue, err error) {
	if db.withLock {
		db.RLock()
		defer db.RUnlock()
	}

	rawVal, ok := db.data[key]
	if !ok {
		return value, ErrKeyNotExists
	}

	value, ok = rawVal.val.(DictValue)
	if !ok {
		return value, ErrKeyTypeMissmatch
	}

	return value, err
}

// SetDict - create/update dict value by key
func (db *DB) SetDict(key string, value DictValue) {
	if db.withLock {
		db.Lock()
		defer db.Unlock()
	}

	item := db.data[key]
	item.val = value
	db.data[key] = item

	return
}

// UpdateDict - update list value by exists key
func (db *DB) UpdateDict(key string, value DictValue) (err error) {
	if db.withLock {
		db.Lock()
		defer db.Unlock()
	}

	if _, ok := db.data[key]; !ok {
		return ErrKeyNotExists
	}

	item := db.data[key]
	item.val = value
	db.data[key] = item

	return nil
}
