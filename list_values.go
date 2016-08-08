package raphanus

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
		return value, ErrKeyNotExists
	}

	value, ok = rawVal.val.(ListValue)
	if !ok {
		return value, ErrKeyTypeMissmatch
	}

	return value, err
}

// SetList - create/update list value by key
func (db *DB) SetList(key string, value ListValue) {
	if db.withLock {
		db.Lock()
		defer db.Unlock()
	}

	item := db.data[key]
	item.val = value
	db.data[key] = item

	return
}

// UpdateList - update list value by exists key
func (db *DB) UpdateList(key string, value ListValue) (err error) {
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

// GetListItem - get one item from list value by key
func (db *DB) GetListItem(key string, index int) (string, error) {
	if db.withLock {
		db.RLock()
		defer db.RUnlock()
	}

	rawVal, ok := db.data[key]
	if !ok {
		return "", ErrKeyNotExists
	}

	valueList, ok := rawVal.val.(ListValue)
	if !ok {
		return "", ErrKeyTypeMissmatch
	}

	if index < 0 || index >= len(valueList) {
		return "", ErrListOutOfRange
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
		return ErrKeyNotExists
	}

	valueList, ok := rawVal.val.(ListValue)
	if !ok {
		return ErrKeyTypeMissmatch
	}

	if index < 0 || index >= len(valueList) {
		return ErrListOutOfRange
	}

	db.data[key].val.(ListValue)[index] = value
	return nil
}
