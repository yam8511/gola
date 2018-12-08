package router

import (
	"gola/app/handler"

	"github.com/gin-gonic/gin"
)

// LoadWebRouter 載入 router 設定
func LoadWebRouter(r *gin.Engine) {
	// 頁面進入點
	r.Use(handler.HomeFile())
}
