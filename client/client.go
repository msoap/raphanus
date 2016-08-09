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

// Cfg - config for New()
type Cfg struct {
	Address string
}

// New - get new client
func New(configs ...Cfg) Client {
	address := defaultAddress
	if len(configs) == 1 {
		address = configs[0].Address
	}

	return Client{
		address: address,
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
	body, err := httpGet(cli.address + APIVersion + "/keys")
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
	body, err := httpDelete(cli.address + APIVersion + "/remove/" + url.QueryEscape(key))
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
	body, err := httpGet(cli.address + APIVersion + "/length")
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

// Integer methods ------------------------------

// GetInt - get int value by key
func (cli Client) GetInt(key string) (int64, error) {
	body, err := httpGet(cli.address + APIVersion + "/int/" + url.QueryEscape(key))
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
func (cli Client) SetInt(key string, value int64, ttl int) (err error) {
	ttlParam := ""
	if ttl > 0 {
		ttlParam = "?ttl=" + url.QueryEscape(strconv.Itoa(ttl))
	}
	postData := []byte(strconv.FormatInt(value, 10))
	body, err := httpPost(cli.address+APIVersion+"/int/"+url.QueryEscape(key)+ttlParam, postData)
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
	body, err := httpPut(cli.address+APIVersion+"/int/"+url.QueryEscape(key), postData)
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
	body, err := httpPost(cli.address+APIVersion+"/int/incr/"+url.QueryEscape(key), nil)
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
	body, err := httpPost(cli.address+APIVersion+"/int/decr/"+url.QueryEscape(key), nil)
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

// String methods ------------------------------

// GetStr - get string value by key
func (cli Client) GetStr(key string) (string, error) {
	body, err := httpGet(cli.address + APIVersion + "/str/" + url.QueryEscape(key))
	if err != nil {
		return "", err
	}

	defer func() {
		if errClose := httpFinalize(body); errClose != nil {
			err = errClose
		}
	}()

	resultRaw := raphanuscommon.OutputGetStr{}
	err = json.NewDecoder(body).Decode(&resultRaw)
	if err != nil {
		return "", err
	}
	if resultRaw.ErrorCode != 0 {
		return "", fmt.Errorf(resultRaw.ErrorMessage)
	}

	return resultRaw.ValueStr, err
}

// SetStr - set string value by key
func (cli Client) SetStr(key string, value string, ttl int) (err error) {
	ttlParam := ""
	if ttl > 0 {
		ttlParam = "?ttl=" + url.QueryEscape(strconv.Itoa(ttl))
	}
	body, err := httpPost(cli.address+APIVersion+"/str/"+url.QueryEscape(key)+ttlParam, []byte(value))
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

// UpdateStr - update string value by key
func (cli Client) UpdateStr(key string, value string) (err error) {
	body, err := httpPut(cli.address+APIVersion+"/str/"+url.QueryEscape(key), []byte(value))
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
