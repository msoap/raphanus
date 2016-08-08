package raphanus

// GetStr - get string value by key
func (db *DB) GetStr(key string) (value string, err error) {
	if db.withLock {
		db.RLock()
		defer db.RUnlock()
	}

	rawVal, ok := db.data[key]
	if !ok {
		return value, ErrKeyNotExists
	}

	value, ok = rawVal.val.(string)
	if !ok {
		return value, ErrKeyTypeMissmatch
	}

	return value, err
}

// SetStr - create/update string value by key
func (db *DB) SetStr(key, value string) {
	if db.withLock {
		db.Lock()
		defer db.Unlock()
	}

	item := db.data[key]
	item.val = value
	db.data[key] = item

	return
}

// UpdateStr - update string value by exists key
func (db *DB) UpdateStr(key, value string) (err error) {
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
