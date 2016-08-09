package raphanusclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"

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

// checkCommonError - check and parse common error from server:
// {"error_code": 0}
// {"error_code":1, "error_message": "..."}
func checkCommonError(body io.ReadCloser) error {
	resultRaw := raphanuscommon.OutputCommon{}
	err := json.NewDecoder(body).Decode(&resultRaw)
	if err != nil {
		return err
	}

	if resultRaw.ErrorCode != 0 {
		return fmt.Errorf(resultRaw.ErrorMessage)
	}

	return nil
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

// Remove - remove key from cache
func (cli Client) Remove(key string) (err error) {
	body, err := httpDelete(defaultAddress + APIVersion + "/remove/" + url.QueryEscape(key))
	if err != nil {
		return err
	}

	defer func() {
		if errClose := httpFinalize(body); errClose != nil {
			err = errClose
		}
	}()

	return checkCommonError(body)
}
