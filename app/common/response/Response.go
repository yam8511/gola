package response

import (
	"gola/app/common/def"
	errorcode "gola/app/common/errorcode"
	"net/http"

	"github.com/gin-gonic/gin"
)

// WithData 回傳附帶資料的選項
func WithData(data interface{}) def.ArgsF {
	return func(opt *def.Option) {
		opt.Data = data
	}
}

// WithText 回傳客製化文字的選項
func WithText(text string) def.ArgsF {
	return func(opt *def.Option) {
		opt.ErrorText = text
	}
}

// WithStatusCode 回傳客製化狀態的選項
func WithStatusCode(status int) def.ArgsF {
	return func(opt *def.Option) {
		opt.HTTPCode = status
	}
}

// Response 回傳
func Response(
	c *gin.Context,
	apiErr errorcode.Error,
	data interface{},
	args ...def.ArgsF,
) interface{} {

	opt := def.WithNewOption(args...)

	var code = 0
	var text string
	if apiErr != nil {
		code = int(apiErr.ErrorCode())
		text = apiErr.ErrorText()
	}

	if opt.ErrorCode != 0 {
		code = opt.ErrorCode
	}

	if opt.ErrorText != "" {
		text = opt.ErrorText
	}

	if opt.Data != nil {
		data = opt.Data
	}

	return API{
		ErrorCode: code,
		ErrorText: text,
		Result:    data,
	}
}

// Success 成功
func Success(c *gin.Context, data interface{}, args ...def.ArgsF) {
	opt := def.Option{
		HTTPCode: http.StatusOK,
	}
	def.WithOption(&opt, args...)

	c.JSON(opt.HTTPCode, Response(c, nil, data, args...))
}

// Failed 失敗
func Failed(c *gin.Context, apiErr errorcode.Error, args ...def.ArgsF) {
	opt := def.Option{
		HTTPCode: http.StatusOK,
	}
	def.WithOption(&opt, args...)

	c.JSON(opt.HTTPCode, Response(c, apiErr, nil, args...))
}

// FailedNow 失敗
func FailedNow(c *gin.Context, apiErr errorcode.Error, args ...def.ArgsF) {
	opt := def.Option{
		HTTPCode: http.StatusOK,
	}
	def.WithOption(&opt, args...)

	c.AbortWithStatusJSON(opt.HTTPCode, Response(c, apiErr, nil, args...))
}
