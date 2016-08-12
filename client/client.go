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
	address  string
	user     string
	password string
}

// Cfg - config for New()
type Cfg struct {
	Address  string
	User     string
	Password string
}

// New - get new client
func New(configs ...Cfg) Client {
	cli := Client{address: defaultAddress}

	if len(configs) == 1 {
		if len(configs[0].Address) > 0 {
			cli.address = configs[0].Address
		}
		if len(configs[0].Password) > 0 {
			cli.user = configs[0].User
			cli.password = configs[0].Password
		}
	}

	return cli
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

// Stat - get some stat from server: version, memory, calls count, etc
func (cli Client) Stat() (result raphanuscommon.Stat, err error) {
	body, err := cli.httpGet(cli.address + APIVersion + "/stat")
	if err != nil {
		return result, err
	}

	defer func() {
		if errClose := httpFinalize(body); errClose != nil {
			err = errClose
		}
	}()

	resultRaw := raphanuscommon.OutputStat{}
	err = json.NewDecoder(body).Decode(&resultRaw)
	if err != nil {
		return result, err
	}
	if resultRaw.ErrorCode != 0 {
		return result, fmt.Errorf(resultRaw.ErrorMessage)
	}

	return resultRaw.Stat, err
}

// Keys - get all keys from cache (response may be too large)
func (cli Client) Keys() (result []string, err error) {
	body, err := cli.httpGet(cli.address + APIVersion + "/keys")
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
	body, err := cli.httpDelete(cli.address + APIVersion + "/remove/" + url.QueryEscape(key))
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
	body, err := cli.httpGet(cli.address + APIVersion + "/length")
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
	body, err := cli.httpGet(cli.address + APIVersion + "/int/" + url.QueryEscape(key))
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

	return resultRaw.ValueInt, nil
}

// SetInt - set int value by key
func (cli Client) SetInt(key string, value int64, ttl int) (err error) {
	ttlParam := ""
	if ttl > 0 {
		ttlParam = "?ttl=" + url.QueryEscape(strconv.Itoa(ttl))
	}
	postData := []byte(strconv.FormatInt(value, 10))
	body, err := cli.httpPost(cli.address+APIVersion+"/int/"+url.QueryEscape(key)+ttlParam, postData)
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
	body, err := cli.httpPut(cli.address+APIVersion+"/int/"+url.QueryEscape(key), postData)
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
	body, err := cli.httpPost(cli.address+APIVersion+"/int/incr/"+url.QueryEscape(key), nil)
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
	body, err := cli.httpPost(cli.address+APIVersion+"/int/decr/"+url.QueryEscape(key), nil)
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
	body, err := cli.httpGet(cli.address + APIVersion + "/str/" + url.QueryEscape(key))
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

	return resultRaw.ValueStr, nil
}

// SetStr - set string value by key
func (cli Client) SetStr(key string, value string, ttl int) (err error) {
	ttlParam := ""
	if ttl > 0 {
		ttlParam = "?ttl=" + url.QueryEscape(strconv.Itoa(ttl))
	}
	body, err := cli.httpPost(cli.address+APIVersion+"/str/"+url.QueryEscape(key)+ttlParam, []byte(value))
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
	body, err := cli.httpPut(cli.address+APIVersion+"/str/"+url.QueryEscape(key), []byte(value))
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

// List methods ------------------------------

// GetList - get list value by key
func (cli Client) GetList(key string) (result raphanuscommon.ListValue, err error) {
	body, err := cli.httpGet(cli.address + APIVersion + "/list/" + url.QueryEscape(key))
	if err != nil {
		return result, err
	}

	defer func() {
		if errClose := httpFinalize(body); errClose != nil {
			err = errClose
		}
	}()

	resultRaw := raphanuscommon.OutputGetList{}
	err = json.NewDecoder(body).Decode(&resultRaw)
	if err != nil {
		return result, err
	}
	if resultRaw.ErrorCode != 0 {
		return result, fmt.Errorf(resultRaw.ErrorMessage)
	}

	return resultRaw.ValueList, nil
}

// SetList - set list value by key
func (cli Client) SetList(key string, value raphanuscommon.ListValue, ttl int) (err error) {
	ttlParam := ""
	if ttl > 0 {
		ttlParam = "?ttl=" + url.QueryEscape(strconv.Itoa(ttl))
	}

	valueJSON, err := json.Marshal(value)
	if err != nil {
		return err
	}

	body, err := cli.httpPost(cli.address+APIVersion+"/list/"+url.QueryEscape(key)+ttlParam, valueJSON)
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

// UpdateList - update list value by key
func (cli Client) UpdateList(key string, value raphanuscommon.ListValue) (err error) {
	valueJSON, err := json.Marshal(value)
	if err != nil {
		return err
	}

	body, err := cli.httpPut(cli.address+APIVersion+"/list/"+url.QueryEscape(key), valueJSON)
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

// GetListItem - get item in list value by key and index
func (cli Client) GetListItem(key string, idx int) (result string, err error) {
	idxParam := "?idx=" + url.QueryEscape(strconv.Itoa(idx))
	body, err := cli.httpGet(cli.address + APIVersion + "/list/item/" + url.QueryEscape(key) + idxParam)
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

	return resultRaw.ValueStr, nil
}

// SetListItem - set item in list value by key and index
func (cli Client) SetListItem(key string, idx int, value string) (err error) {
	idxParam := "?idx=" + url.QueryEscape(strconv.Itoa(idx))
	body, err := cli.httpPut(cli.address+APIVersion+"/list/item/"+url.QueryEscape(key)+idxParam, []byte(value))
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

// Dict methods ------------------------------

// GetDict - get dict value by key
func (cli Client) GetDict(key string) (result raphanuscommon.DictValue, err error) {
	body, err := cli.httpGet(cli.address + APIVersion + "/dict/" + url.QueryEscape(key))
	if err != nil {
		return result, err
	}

	defer func() {
		if errClose := httpFinalize(body); errClose != nil {
			err = errClose
		}
	}()

	resultRaw := raphanuscommon.OutputGetDict{}
	err = json.NewDecoder(body).Decode(&resultRaw)
	if err != nil {
		return result, err
	}
	if resultRaw.ErrorCode != 0 {
		return result, fmt.Errorf(resultRaw.ErrorMessage)
	}

	return resultRaw.ValueDict, nil
}

// SetDict - set dict value by key
func (cli Client) SetDict(key string, value raphanuscommon.DictValue, ttl int) (err error) {
	ttlParam := ""
	if ttl > 0 {
		ttlParam = "?ttl=" + url.QueryEscape(strconv.Itoa(ttl))
	}

	valueJSON, err := json.Marshal(value)
	if err != nil {
		return err
	}

	body, err := cli.httpPost(cli.address+APIVersion+"/dict/"+url.QueryEscape(key)+ttlParam, valueJSON)
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

// UpdateDict - update dict value by key
func (cli Client) UpdateDict(key string, value raphanuscommon.DictValue) (err error) {
	valueJSON, err := json.Marshal(value)
	if err != nil {
		return err
	}

	body, err := cli.httpPut(cli.address+APIVersion+"/dict/"+url.QueryEscape(key), valueJSON)
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

// GetDictItem - get item form dict value by key and dict key
func (cli Client) GetDictItem(key string, dkey string) (result string, err error) {
	dkeyParam := "?dkey=" + url.QueryEscape(dkey)
	body, err := cli.httpGet(cli.address + APIVersion + "/dict/item/" + url.QueryEscape(key) + dkeyParam)
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

	return resultRaw.ValueStr, nil
}

// SetDictItem - set item in dict value by key and dict key
func (cli Client) SetDictItem(key string, dkey string, value string) (err error) {
	dkeyParam := "?dkey=" + url.QueryEscape(dkey)
	body, err := cli.httpPut(cli.address+APIVersion+"/dict/item/"+url.QueryEscape(key)+dkeyParam, []byte(value))
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

// RemoveDictItem - remove one item from dict value by key and dict key
func (cli Client) RemoveDictItem(key string, dkey string) (err error) {
	dkeyParam := "?dkey=" + url.QueryEscape(dkey)
	body, err := cli.httpDelete(cli.address + APIVersion + "/dict/item/" + url.QueryEscape(key) + dkeyParam)
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
