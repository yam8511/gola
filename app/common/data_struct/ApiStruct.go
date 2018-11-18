package datastruct

// API 格式
type API struct {
	ErrorCode int         `json:"error_code"`
	ErrorText string      `json:"error_text"`
	Result    interface{} `json:"result"`
}
