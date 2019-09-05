package router

import (
	"gola/app/handler"
	"gola/app/middleware"

	"github.com/gin-gonic/gin"
)

// LoadAPIRouter 載入 router 設定
func LoadAPIRouter(r *gin.Engine) {
	api := r.Group(
		"/api",
		middleware.GetMiddleware("check_google_login"),
	)
	api.GET("/hello", handler.API)
	api.POST("/seed", handler.Seed)
}
