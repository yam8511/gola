package router

import (
	"gola/app/handler"
	"net/url"

	"github.com/gin-gonic/gin"
)

// LoadWebRouter 載入 router 設定
func LoadWebRouter(r *gin.Engine) {
	r.GET("/favicon.png", func(c *gin.Context) {
		referer, err := url.Parse(c.Request.Referer())
		if err == nil && referer.Path == "/star/" {
			c.File("./public/star/favicon.png")
			return
		}

		c.Status(404)
	})
	// 頁面進入點
	r.NoRoute(handler.HomeFile())
	// r.Use(handler.HomeFile())
}
