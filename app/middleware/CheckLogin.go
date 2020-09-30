package middleware

import (
	"gola/app/common/errorcode"
	"gola/app/common/response"
	google "gola/app/service/google"

	"github.com/gin-gonic/gin"
)

// 驗證舊系統登入
func checkGoogleLogin(c *gin.Context) {
	sid, err := c.Cookie("sid")
	if err != nil {
		apiErr := errorcode.Code_No_Cookie.New("")
		response.FailedNow(c, apiErr)
		return
	}

	isLogin, apiErr := google.CheckLogin(sid)
	if apiErr != nil {
		response.FailedNow(c, apiErr)
		return
	}

	if !isLogin {
		apiErr := errorcode.Code_No_Login.New("")
		response.FailedNow(c, apiErr)
	}

	return
}
