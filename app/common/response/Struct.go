package response

import "gola/app/common/errorcode"

// API 格式
type API struct {
	ErrorCode errorcode.Code `json:"error_code"`
	ErrorText string         `json:"error_text"`
	Result    interface{}    `json:"result"`
}
