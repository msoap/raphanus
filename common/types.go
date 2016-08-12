package raphanuscommon

// ListValue - list value type
type ListValue []string

// DictValue - dict value type
type DictValue map[string]string

// Stat - some stat: version, memory, calls count, etc
type Stat struct {
	Version string `json:"version"`
}

// JSON types for responses from server

// OutputCommon - common part of all responses
type OutputCommon struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// OutputStat - output for /stat
type OutputStat struct {
	OutputCommon
	Stat
}

// OutputKeys - output for /keys
type OutputKeys struct {
	OutputCommon
	Keys []string `json:"keys"`
}

// OutputLength - output for /length
type OutputLength struct {
	OutputCommon
	Length int `json:"length"`
}

// OutputGetInt - output for /int/:key
type OutputGetInt struct {
	OutputCommon
	ValueInt int64 `json:"value_int"`
}

// OutputGetStr - output for /str/:key
type OutputGetStr struct {
	OutputCommon
	ValueStr string `json:"value_str"`
}

// OutputGetList - output for /list/:key
type OutputGetList struct {
	OutputCommon
	ValueList ListValue `json:"value_list"`
}
