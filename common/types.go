package raphanuscommon

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
