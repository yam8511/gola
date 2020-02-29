package request

import (
	"gola/app/common/constant"

	"github.com/gin-gonic/gin"
)

// GetLang 取語系
func GetLang(c *gin.Context) string {
	lang, ok := c.GetQuery(constant.LangHeaderKey)
	if !ok {
		lang = c.GetHeader(constant.LangHeaderKey)
	}
	switch lang {
	case "en":
		lang = "en"
	case "zh-tw", "tw":
		lang = "zh-tw"
	default:
		lang = "zh-cn"
	}
	return lang
}
