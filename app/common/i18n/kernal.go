package i18n

import (
	"gola/app/common/request"

	"github.com/gin-gonic/gin"
)

// Trans 取得翻譯
func Trans(code, lang string) (value string) {
	switch lang {
	case "en":
		value = enBook.Trans(code)
	case "zh-cn", "cn":
		value = cnBook.Trans(code)
	default: // zh-cn
		value = twBook.Trans(code)
	}

	return
}

// TransContext 取連線語系
func TransContext(c *gin.Context, code string) (value string) {
	lang := request.GetLang(c)
	return Trans(code, lang)
}
