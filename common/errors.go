package raphanuscommon

// RaphError - common error type
type RaphError struct {
	Code    int
	Message string
}

// Error - implement error interface
func (err RaphError) Error() string {
	return err.Message
}

// Common errors
var (
	ErrBadRequest       = RaphError{1, "Bad request"}
	ErrKeyNotExists     = RaphError{2, "Key not exists"}
	ErrKeyTypeMissmatch = RaphError{3, "The type does not match"}
	ErrListOutOfRange   = RaphError{4, "List index is out of range"}
	ErrDictKeyNotExists = RaphError{5, "Dict, key not exists"}
	ErrDictKeyIsEmpty   = RaphError{6, "Key or dict key is empty"}
	ErrTTLIsntCorrect   = RaphError{7, "TTL parameter isn't correct"}
	ErrKeyIsNotValid    = RaphError{8, "Key is not valid"}
)
