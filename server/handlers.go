package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/msoap/raphanus"
)

// outputCommon - common part of all responses
type outputCommon struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// allocate result for success calls
var outputCommonOK = outputCommon{}

/*
handlerKeys - get all keys

curl http://localhost:8771/v1/keys
result:
	{"error_code":0,"keys":["k1", "k2"]}
*/
func (app *server) handlerKeys(ctx echo.Context) error {
	type outputKeys struct {
		outputCommon
		Keys []string `json:"keys"`
	}

	result := outputKeys{Keys: app.raphanus.Keys()}
	return ctx.JSON(http.StatusOK, result)
}

/*
handlerRemoveKey - remove key

curl -s -X DELETE http://localhost:8771/v1/remove/k1
result:
	{"error_code":0}
*/
func (app *server) handlerRemoveKey(ctx echo.Context) error {
	key := ctx.Param("key")
	err := app.raphanus.Remove(key)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, outputCommon{ErrorCode: 1, ErrorMessage: err.Error()})
	}

	return ctx.JSON(http.StatusOK, outputCommonOK)
}

/*
handlerLength - get count of keys

curl -s http://localhost:8771/v1/length
result:
	{"error_code":0, "length": 3}
*/
func (app *server) handlerLength(ctx echo.Context) error {
	type outputLength struct {
		outputCommon
		Length int `json:"length"`
	}
	len := app.raphanus.Len()
	return ctx.JSON(http.StatusOK, outputLength{Length: len})
}

// Integer methods ------------------------------

type outputGetInt struct {
	outputCommon
	ValueInt int64 `json:"value_int"`
}

/*
getInt - get one integer value by key

curl -s http://localhost:8771/v1/int/k1
result:
	{"error_code":0,"value_int":737}
*/
func (app *server) getInt(ctx echo.Context) error {
	key := ctx.Param("key")
	valueInt, err := app.raphanus.GetInt(key)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, outputCommon{ErrorCode: 1, ErrorMessage: err.Error()})
	}

	result := outputGetInt{ValueInt: valueInt}
	return ctx.JSON(http.StatusOK, result)
}

/*
setInt - set one integer value by key

curl -s -X POST -d 123 http://localhost:8771/v1/int/k1
result:
	{"error_code":0}
*/
func (app *server) setInt(ctx echo.Context) error {
	newIntValue, err := getBodyAsInt64(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, outputCommon{ErrorCode: 1, ErrorMessage: err.Error()})
	}

	key := ctx.Param("key")
	app.raphanus.SetInt(key, newIntValue)

	return ctx.JSON(http.StatusOK, outputCommonOK)
}

/*
updateInt - set one integer value by key

curl -s -X PUT -d 127 http://localhost:8771/v1/int/k1
result:
	{"error_code":0}
*/
func (app *server) updateInt(ctx echo.Context) error {
	newIntValue, err := getBodyAsInt64(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, outputCommon{ErrorCode: 1, ErrorMessage: err.Error()})
	}

	key := ctx.Param("key")
	if err := app.raphanus.UpdateInt(key, newIntValue); err != nil {
		return ctx.JSON(http.StatusBadRequest, outputCommon{ErrorCode: 1, ErrorMessage: err.Error()})
	}

	return ctx.JSON(http.StatusOK, outputCommonOK)
}

/*
incrInt - increment one value

curl -s -X POST http://localhost:8771/v1/int/incr/k1
result:
	{"error_code":0,"value_int":738}
*/
func (app *server) incrInt(ctx echo.Context) error {
	key := ctx.Param("key")

	var (
		err      error
		valueInt int64
	)
	app.raphanus.UnderLock(func() {
		err = app.raphanus.IncrInt(key)
		if err != nil {
			return
		}
		valueInt, err = app.raphanus.GetInt(key)
		if err != nil {
			return
		}
	})

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, outputCommon{ErrorCode: 1, ErrorMessage: err.Error()})
	}

	result := outputGetInt{ValueInt: valueInt}
	return ctx.JSON(http.StatusOK, result)
}

/*
decrInt - decrement one value

curl -s -X POST http://localhost:8771/v1/int/decr/k1
result:
	{"error_code":0,"value_int":736}
*/
func (app *server) decrInt(ctx echo.Context) error {
	key := ctx.Param("key")

	var (
		err      error
		valueInt int64
	)
	app.raphanus.UnderLock(func() {
		err = app.raphanus.DecrInt(key)
		if err != nil {
			return
		}
		valueInt, err = app.raphanus.GetInt(key)
		if err != nil {
			return
		}
	})

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, outputCommon{ErrorCode: 1, ErrorMessage: err.Error()})
	}

	result := outputGetInt{ValueInt: valueInt}
	return ctx.JSON(http.StatusOK, result)
}

// getBodyAsInt64 - get body of request as int64
func getBodyAsInt64(ctx echo.Context) (int64, error) {
	// read first 20 bytes
	limitBody := io.LimitReader(ctx.Request().Body(), 20)
	bytes, err := ioutil.ReadAll(limitBody)
	if err != nil {
		return 0, err
	}

	intValue, err := strconv.ParseInt(string(bytes), 10, 64)
	if err != nil {
		return 0, err
	}
	return intValue, nil
}

// String methods ------------------------------

type outputGetStr struct {
	outputCommon
	ValueStr string `json:"value_str"`
}

/*
getStr - get one string value by key

curl -s http://localhost:8771/v1/str/k1
result:
	{"error_code":0,"value_str":"string value"}
*/
func (app *server) getStr(ctx echo.Context) error {
	key := ctx.Param("key")
	valueStr, err := app.raphanus.GetStr(key)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, outputCommon{ErrorCode: 1, ErrorMessage: err.Error()})
	}

	result := outputGetStr{ValueStr: valueStr}
	return ctx.JSON(http.StatusOK, result)
}

/*
setStr - set one string value by key

curl -s -X POST -d "string value" http://localhost:8771/v1/str/k1
result:
	{"error_code":0}
*/
func (app *server) setStr(ctx echo.Context) error {
	newStrValue, err := getBodyAsString(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, outputCommon{ErrorCode: 1, ErrorMessage: err.Error()})
	}

	key := ctx.Param("key")
	app.raphanus.SetStr(key, newStrValue)

	return ctx.JSON(http.StatusOK, outputCommonOK)
}

/*
updateStr - set one string value by key

curl -s -X PUT -d "new value" http://localhost:8771/v1/str/k1
result:
	{"error_code":0}
*/
func (app *server) updateStr(ctx echo.Context) error {
	newStrValue, err := getBodyAsString(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, outputCommon{ErrorCode: 1, ErrorMessage: err.Error()})
	}

	key := ctx.Param("key")
	if err := app.raphanus.UpdateStr(key, newStrValue); err != nil {
		return ctx.JSON(http.StatusBadRequest, outputCommon{ErrorCode: 1, ErrorMessage: err.Error()})
	}

	return ctx.JSON(http.StatusOK, outputCommonOK)
}

// getBodyAsString - get body of request as string
func getBodyAsString(ctx echo.Context) (string, error) {
	limitBody := io.LimitReader(ctx.Request().Body(), raphanus.MaxStringValueLength)
	bytes, err := ioutil.ReadAll(limitBody)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// List methods ------------------------------

/*
getList - get one list value by key

curl -s http://localhost:8771/v1/list/k1
result:
	{"error_code":0,"value_list":["l1", "l2", "l3"]}
*/
func (app *server) getList(ctx echo.Context) error {
	type outputGetList struct {
		outputCommon
		ValueList []string `json:"value_list"`
	}

	key := ctx.Param("key")
	valueList, err := app.raphanus.GetList(key)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, outputCommon{ErrorCode: 1, ErrorMessage: err.Error()})
	}

	result := outputGetList{ValueList: valueList}
	return ctx.JSON(http.StatusOK, result)
}

/*
setList - set one list value by key

curl -s -X POST -H 'Content-Type: application/json' -d '["l1", "l2", "l3"]' http://localhost:8771/v1/list/k1
result:
	{"error_code":0}
*/
func (app *server) setList(ctx echo.Context) error {
	newListValue := []string{}
	if err := ctx.Bind(&newListValue); err != nil {
		return ctx.JSON(http.StatusBadRequest, outputCommon{ErrorCode: 1, ErrorMessage: err.Error()})
	}

	key := ctx.Param("key")
	app.raphanus.SetList(key, newListValue)

	return ctx.JSON(http.StatusOK, outputCommonOK)
}

/*
updateList - update one list value by exists key

curl -s -X PUT -H 'Content-Type: application/json' -d '["l1", "l2"]' http://localhost:8771/v1/list/k1
result:
	{"error_code":0}
*/
func (app *server) updateList(ctx echo.Context) error {
	newListValue := []string{}
	if err := ctx.Bind(&newListValue); err != nil {
		return ctx.JSON(http.StatusBadRequest, outputCommon{ErrorCode: 1, ErrorMessage: err.Error()})
	}

	key := ctx.Param("key")
	if err := app.raphanus.UpdateList(key, newListValue); err != nil {
		return ctx.JSON(http.StatusBadRequest, outputCommon{ErrorCode: 1, ErrorMessage: err.Error()})
	}

	return ctx.JSON(http.StatusOK, outputCommonOK)
}

/*
getListItem - get one item from list value by key and index in list

curl -s 'http://localhost:8771/v1/list/item/k1?idx=1'
result:
	{"result": "ok", "value_str": "l2"}
*/
func (app *server) getListItem(ctx echo.Context) error {
	index, err := strconv.Atoi(ctx.QueryParam("idx"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, outputCommon{ErrorCode: 1, ErrorMessage: err.Error()})
	}

	key := ctx.Param("key")
	valueStr, err := app.raphanus.GetListItem(key, index)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, outputCommon{ErrorCode: 1, ErrorMessage: err.Error()})
	}

	return ctx.JSON(http.StatusOK, outputGetStr{ValueStr: valueStr})
}

/*
setListItem - set one item on list value by key and index in list

curl -s -X PUT -d "l3" 'http://localhost:8771/v1/list/item/k1?idx=1'
result:
	{"error_code":0}
*/
func (app *server) setListItem(ctx echo.Context) error {
	index, err := strconv.Atoi(ctx.QueryParam("idx"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, outputCommon{ErrorCode: 1, ErrorMessage: err.Error()})
	}

	newStrValue, err := getBodyAsString(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, outputCommon{ErrorCode: 1, ErrorMessage: err.Error()})
	}

	key := ctx.Param("key")
	if err := app.raphanus.SetListItem(key, index, newStrValue); err != nil {
		return ctx.JSON(http.StatusBadRequest, outputCommon{ErrorCode: 1, ErrorMessage: err.Error()})
	}

	return ctx.JSON(http.StatusOK, outputCommonOK)
}

// Dict methods ------------------------------

/*
getDict - get one dict value by key

curl -s http://localhost:8771/v1/dict/k1
result:
	{"error_code":0, "value_list": {"dk1": "v1", "dk2": "v2"}}
*/
func (app *server) getDict(ctx echo.Context) error {
	type outputGetDict struct {
		outputCommon
		ValueDict raphanus.DictValue `json:"value_dict"`
	}

	key := ctx.Param("key")
	valueDict, err := app.raphanus.GetDict(key)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, outputCommon{ErrorCode: 1, ErrorMessage: err.Error()})
	}

	result := outputGetDict{ValueDict: valueDict}
	return ctx.JSON(http.StatusOK, result)
}

/*
setDict - set one dict value by key

curl -s -X POST -H 'Content-Type: application/json' -d '{"dk1": "v1", "dk2": "v2"}' http://localhost:8771/v1/dict/k1
result:
	{"error_code":0}
*/
func (app *server) setDict(ctx echo.Context) error {
	newDictValue := raphanus.DictValue{}
	if err := ctx.Bind(&newDictValue); err != nil {
		return ctx.JSON(http.StatusBadRequest, outputCommon{ErrorCode: 1, ErrorMessage: err.Error()})
	}

	key := ctx.Param("key")
	app.raphanus.SetDict(key, newDictValue)

	return ctx.JSON(http.StatusOK, outputCommonOK)
}

/*
updateDict - update one dict value by exists key

curl -s -X PUT -H 'Content-Type: application/json' -d '{"dk1": "v33"}' http://localhost:8771/v1/dict/k1
result:
	{"error_code":0}
*/
func (app *server) updateDict(ctx echo.Context) error {
	newDictValue := raphanus.DictValue{}
	if err := ctx.Bind(&newDictValue); err != nil {
		return ctx.JSON(http.StatusBadRequest, outputCommon{ErrorCode: 1, ErrorMessage: err.Error()})
	}

	key := ctx.Param("key")
	if err := app.raphanus.UpdateDict(key, newDictValue); err != nil {
		return ctx.JSON(http.StatusBadRequest, outputCommon{ErrorCode: 1, ErrorMessage: err.Error()})
	}

	return ctx.JSON(http.StatusOK, outputCommonOK)
}
