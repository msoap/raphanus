package raphanus

// GetInt - get integer value by key
func (db *DB) GetInt(key string) (int64, error) {
	if db.withLock {
		db.RLock()
		defer db.RUnlock()
	}

	rawVal, ok := db.data[key]
	if !ok {
		return 0, ErrKeyNotExists
	}

	value, ok := rawVal.val.(int64)
	if !ok {
		return 0, ErrKeyTypeMissmatch
	}

	return value, nil
}

// SetInt - create/update integer value by key
func (db *DB) SetInt(key string, value int64) {
	if db.withLock {
		db.Lock()
		defer db.Unlock()
	}

	item := db.data[key]
	item.val = value
	db.data[key] = item

	return
}

// UpdateInt - update integer value by exists key
func (db *DB) UpdateInt(key string, value int64) (err error) {
	if db.withLock {
		db.Lock()
		defer db.Unlock()
	}

	_, ok := db.data[key]
	if !ok {
		return ErrKeyNotExists
	}

	item := db.data[key]
	item.val = value
	db.data[key] = item

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
		return ErrKeyNotExists
	}

	_, ok = db.data[key].val.(int64)
	if !ok {
		return ErrKeyTypeMissmatch
	}

	item := db.data[key]
	item.val = item.val.(int64) + value
	db.data[key] = item

	return nil
}
