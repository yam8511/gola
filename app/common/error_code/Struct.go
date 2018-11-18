package errorcode

import "fmt"

// APIError API錯誤格式
type APIError struct {
	ErrorCode int    `json:"error_code"`
	ErrorText string `json:"error_text"`
}

// Error API錯誤訊息
func (e APIError) Error() string {
	return fmt.Sprintf("%s (%d)", e.ErrorText, e.ErrorCode)
}
