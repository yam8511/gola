package errorcode

import (
	"strconv"
)

// APIError API錯誤格式
type APIError struct {
	Code    int    `json:"error_code"`
	Text    string `json:"error_text"`
	showLog bool
}

// ErrorCode 錯誤代碼
func (e APIError) ErrorCode() int {
	return e.Code
}

// ErrorText 錯誤訊息
func (e APIError) ErrorText() string {
	return e.Text
}

// Error API錯誤訊息
func (e APIError) Error() string {
	return e.Text + " (" + strconv.Itoa(e.Code) + ")"
}

// GetAPIError 由錯誤碼取得API回傳
func GetAPIError(text string, err error) Error {
	return newAPIError(text, err)
}
