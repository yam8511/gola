package middleware

import (
	dataStruct "gola/app/common/data_struct"
	errorCode "gola/app/common/error_code"
	google "gola/app/service/google"
	"gola/internal/bootstrap"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 驗證舊系統登入
func checkGoogleLogin(c *gin.Context) {
	sid, err := c.Cookie("sid")
	if err != nil {
		apiErr := errorCode.GetAPIError(1000)
		c.AbortWithStatusJSON(http.StatusOK, dataStruct.API{
			ErrorCode: apiErr.ErrorCode,
			ErrorText: apiErr.ErrorText,
		})
		return
	}

	isLogin, apiErr := google.CheckLogin(sid)
	if apiErr.ErrorCode != 0 {
		c.AbortWithStatusJSON(http.StatusOK, dataStruct.API{
			ErrorCode: apiErr.ErrorCode,
			ErrorText: apiErr.ErrorText,
		})
		return
	}

	if !isLogin {
		bootstrap.WriteLog("INFO", "登入失敗！")
		apiErr := errorCode.GetAPIError(1004)
		c.AbortWithStatusJSON(http.StatusOK, dataStruct.API{
			ErrorCode: apiErr.ErrorCode,
			ErrorText: apiErr.ErrorText,
		})
	} else {
		bootstrap.WriteLog("INFO", "登入成功！")
	}

	return
}
