package raphanusclient

import (
	"encoding/json"

	"github.com/msoap/raphanus/common"
)

const (
	defaultAddress = "http://" + raphanuscommon.DefaultHost + ":" + raphanuscommon.DefaultPort
	timeout        = 60 // HTTP client timout
	// APIVersion - prefix for path in URL
	APIVersion = "/v1"
)

// Client - client object
type Client struct {
	address string
}

// New - get new client
func New() Client {
	return Client{
		address: defaultAddress,
	}
}

// Keys - get all keys from cache (response may be too large)
func (cli Client) Keys() (result []string, err error) {
	body, err := httpGet(defaultAddress + APIVersion + "/keys")
	if err != nil {
		return result, err
	}

	defer func() {
		if errClose := httpFinalize(body); errClose != nil {
			err = errClose
		}
	}()

	resultRaw := raphanuscommon.OutputKeys{}
	err = json.NewDecoder(body).Decode(&resultRaw)
	if err != nil {
		return result, err
	}

	return resultRaw.Keys, err
}
