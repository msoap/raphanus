package raphanuscommon

// JSON types for responses from server

// OutputCommon - common part of all responses
type OutputCommon struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message,omitempty"`
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
