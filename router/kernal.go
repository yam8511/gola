package router

import (
	"gola/app/middleware"
	"gola/internal/bootstrap"
	defaultR "gola/router/default"

	"github.com/gin-gonic/gin"
)

// RouteProvider 路由提供者
func RouteProvider(r *gin.Engine) {
	// 全域載入
	r.Use(middleware.GlobalMiddlewares()...)

	Conf := bootstrap.GetAppConf()
	switch Conf.App.Site {
	case "admin":
		// 專屬admin的route
	case "member":
		// 專屬member的route
	default:
		defaultR.LoadRoutes(r)
	}
}
