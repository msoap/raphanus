package raphanusclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"

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
func checkCommonError(body io.Reader) error {
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
	if resultRaw.ErrorCode != 0 {
		return result, fmt.Errorf(resultRaw.ErrorMessage)
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

// Length - get count of keys
func (cli Client) Length() (int, error) {
	body, err := httpGet(defaultAddress + APIVersion + "/length")
	if err != nil {
		return 0, err
	}

	defer func() {
		if errClose := httpFinalize(body); errClose != nil {
			err = errClose
		}
	}()

	resultRaw := raphanuscommon.OutputLength{}
	err = json.NewDecoder(body).Decode(&resultRaw)
	if err != nil {
		return 0, err
	}
	if resultRaw.ErrorCode != 0 {
		return 0, fmt.Errorf(resultRaw.ErrorMessage)
	}

	return resultRaw.Length, err
}

// GetInt - get int value by key
func (cli Client) GetInt(key string) (int64, error) {
	body, err := httpGet(defaultAddress + APIVersion + "/int/" + url.QueryEscape(key))
	if err != nil {
		return 0, err
	}

	defer func() {
		if errClose := httpFinalize(body); errClose != nil {
			err = errClose
		}
	}()

	resultRaw := raphanuscommon.OutputGetInt{}
	err = json.NewDecoder(body).Decode(&resultRaw)
	if err != nil {
		return 0, err
	}
	if resultRaw.ErrorCode != 0 {
		return 0, fmt.Errorf(resultRaw.ErrorMessage)
	}

	return resultRaw.ValueInt, err
}

// SetInt - set int value by key
func (cli Client) SetInt(key string, value int64) (err error) {
	postData := []byte(strconv.FormatInt(value, 10))
	body, err := httpPost(defaultAddress+APIVersion+"/int/"+url.QueryEscape(key), postData)
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

// UpdateInt - update int value by key
func (cli Client) UpdateInt(key string, value int64) (err error) {
	postData := []byte(strconv.FormatInt(value, 10))
	body, err := httpPut(defaultAddress+APIVersion+"/int/"+url.QueryEscape(key), postData)
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

// IncrInt - increment int value by key
func (cli Client) IncrInt(key string) (err error) {
	body, err := httpPost(defaultAddress+APIVersion+"/int/incr/"+url.QueryEscape(key), nil)
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

// DecrInt - decrement int value by key
func (cli Client) DecrInt(key string) (err error) {
	body, err := httpPost(defaultAddress+APIVersion+"/int/decr/"+url.QueryEscape(key), nil)
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
