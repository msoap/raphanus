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

	return value, nil
}

// SetDict - create/update dict value by key
func (db *DB) SetDict(key string, value DictValue, ttl int) {
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

// UpdateDict - update dict value by exists key
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

// GetDictItem - get item from dict value by exists key
func (db *DB) GetDictItem(key string, dictKey string) (string, error) {
	if db.withLock {
		db.RLock()
		defer db.RUnlock()
	}

	if err := db.validateDictParams(key, dictKey); err != nil {
		return "", err
	}

	return db.data[key].val.(DictValue)[dictKey], nil
}

// SetDictItem - set item on dict value by exists key
func (db *DB) SetDictItem(key, dictKey, dictValue string) error {
	if db.withLock {
		db.RLock()
		defer db.RUnlock()
	}

	if err := db.validateDictParams(key, dictKey); err != nil {
		return err
	}

	db.data[key].val.(DictValue)[dictKey] = dictValue

	return nil
}

// RemoveDictItem - remove item on dict value by exists key
func (db *DB) RemoveDictItem(key, dictKey string) error {
	if db.withLock {
		db.RLock()
		defer db.RUnlock()
	}

	if err := db.validateDictParams(key, dictKey); err != nil {
		return err
	}

	delete(db.data[key].val.(DictValue), dictKey)

	return nil
}

func (db *DB) validateDictParams(key, dictKey string) error {
	if len(key) == 0 || len(dictKey) == 0 {
		return ErrDictKeyIsEmpty
	}

	rawVal, ok := db.data[key]
	if !ok {
		return ErrKeyNotExists
	}

	value, ok := rawVal.val.(DictValue)
	if !ok {
		return ErrKeyTypeMissmatch
	}

	if _, ok := value[dictKey]; !ok {
		return ErrDictKeyNotExists
	}

	return nil
}
