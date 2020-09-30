package errorcode

import (
	"fmt"
	"gola/internal/logger"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
)

var _ Error = APIError{}

// APIError API錯誤格式
type APIError struct {
	Code Code   `json:"error_code"`
	Text string `json:"error_text"`
}

// ErrorCode 錯誤代碼
func (e APIError) ErrorCode() Code {
	return e.Code
}

// ErrorText 錯誤訊息
func (e APIError) ErrorText() string {
	return e.Text
}

// Error API錯誤訊息
func (e APIError) Error() string {
	return e.Text + " (" + strconv.Itoa(int(e.Code)) + ")"
}

// 錯誤代碼型態
type Code int

func (c Code) New(text string, v ...interface{}) Error {
	if text != "" {
		if c == Code_Panic {
			name, file, line := funcinfo(4)
			if strings.Contains(file, "panic") {
				name, file, line = funcinfo(5)
			} else if strings.Contains(name, "goexit") {
				name, file, line = funcinfo(2)
			}
			logger.Danger(
				"🚑 🚑 🚑 \n%s\n🎃  📦  檔案: %s:%d  🧩  Func: %s  🐞  內部發生`panic`錯誤 [%d], "+text+" 🎃\n",
				append([]interface{}{string(debug.Stack()), file, line, name, c}, v...)...,
			)
		} else {
			name, file, line := funcinfo(2)
			logger.Warn(
				"🎃  📦  檔案: %s:%d  🧩  Func: %s  🐞  內部發生錯誤 [%d], "+text+" 🎃\n",
				append([]interface{}{file, line, name, c}, v...)...,
			)
		}
	}

	return APIError{
		Code: c,
		Text: fmt.Sprintf(text, v...),
	}
}

func funcinfo(skip int) (name, file string, line int) {
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		name = runtime.FuncForPC(pc).Name()
	}
	return name, file, line
}
