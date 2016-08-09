package raphanusclient

import (
	"github.com/msoap/raphanus/common"
)

const (
	defaultAddress = "http://" + raphanuscommon.DefaultHost + ":" + raphanuscommon.DefaultPort
)

// DB - client object
type DB struct {
	address string
}

// New - get new client
func New() DB {
	return DB{
		address: defaultAddress,
	}
}
