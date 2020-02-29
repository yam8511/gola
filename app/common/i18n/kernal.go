package i18n

import (
	"gola/app/common/request"
	"sync"

	"github.com/gin-gonic/gin"
)

type book struct {
	sync.Once
	content []byte
}

// Trans 取得翻譯
func Trans(code, lang string) (value string) {
	switch lang {
	case "en":
		value = transEN(code)
	case "zh-tw", "tw":
		value = transTW(code)
	default: // zh-cn
		value = transCN(code)
	}

	return
}

// TransContext 取連線語系
func TransContext(c *gin.Context, code string) (value string) {
	lang := request.GetLang(c)
	return Trans(code, lang)
}
