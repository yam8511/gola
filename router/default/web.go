package router

import (
	"gola/app/handler"
	"gola/docs"
	"net/url"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
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

	docs.Init()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 頁面進入點
	r.NoRoute(handler.HomeFile())
}
