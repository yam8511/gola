package middleware

import (
	errorCode "gola/app/common/errorcode"
	"gola/app/common/response"
	google "gola/app/service/google"

	"github.com/gin-gonic/gin"
)

// 驗證舊系統登入
func checkGoogleLogin(c *gin.Context) {
	sid, err := c.Cookie("sid")
	if err != nil {
		apiErr := errorCode.GetAPIError("no_cookie", nil)
		response.FailedNow(c, apiErr)
		return
	}

	isLogin, apiErr := google.CheckLogin(sid)
	if apiErr != nil {
		response.FailedNow(c, apiErr)
		return
	}

	if !isLogin {
		apiErr := errorCode.GetAPIError("not_login", nil)
		response.FailedNow(c, apiErr)
	}

	return
}
