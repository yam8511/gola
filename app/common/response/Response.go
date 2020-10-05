package response

import (
	"gola/app/common/errorcode"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 回傳
func Response(
	c *gin.Context,
	apiErr errorcode.Error,
	data interface{},
	args ...ArgsF,
) interface{} {

	opt := WithNewOption(args...)

	var code errorcode.Code
	var text string
	if apiErr != nil {
		code = apiErr.ErrorCode()
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
func Success(c *gin.Context, data interface{}, args ...ArgsF) {
	opt := Option{
		HTTPCode: http.StatusOK,
	}
	WithOption(&opt, args...)

	c.JSON(opt.HTTPCode, Response(c, nil, data, args...))
}

// Failed 失敗
func Failed(c *gin.Context, apiErr errorcode.Error, args ...ArgsF) {
	opt := Option{
		HTTPCode: http.StatusOK,
	}
	WithOption(&opt, args...)

	c.JSON(opt.HTTPCode, Response(c, apiErr, nil, args...))
}

// FailedNow 失敗
func FailedNow(c *gin.Context, apiErr errorcode.Error, args ...ArgsF) {
	opt := Option{
		HTTPCode: http.StatusOK,
	}
	WithOption(&opt, args...)

	c.AbortWithStatusJSON(opt.HTTPCode, Response(c, apiErr, nil, args...))
}
