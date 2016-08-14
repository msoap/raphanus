package raphanuscommon

// ListValue - list value type
type ListValue []string

// DictValue - dict value type
type DictValue map[string]string

// Stat - some stat: version, memory, GC, etc
type Stat struct {
	Version        string `json:"version"`
	MemAlloc       uint64 `json:"mem_alloc"`
	MemTotalAlloc  uint64 `json:"mem_total_alloc"`
	MemMallocs     uint64 `json:"mem_malloc"`
	MemFrees       uint64 `json:"mem_frees"`
	MemHeapObjects uint64 `json:"mem_heap_objects"`
	GCPauseTotalNs uint64 `json:"gc_pause_total_ns"`
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

// OutputGetDict - output for /dict/:key
type OutputGetDict struct {
	OutputCommon
	ValueDict DictValue `json:"value_dict"`
}
