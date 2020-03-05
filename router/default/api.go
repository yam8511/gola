package router

import (
	"gola/app/handler"

	"github.com/gin-gonic/gin"
)

// LoadAPIRouter 載入 router 設定
func LoadAPIRouter(r *gin.Engine) {
	api := r.Group(
		"/api",
		// middleware.Middleware("check_google_login"),
	)
	api.GET("/hello", handler.API)
	api.POST("/config", handler.Config)
	// api.POST("/seed", handler.Seed)
}
