package router

import (
	"gola/app/middleware"
	"gola/internal/bootstrap"
	defaultR "gola/router/default"

	"github.com/gin-gonic/gin"
)

// RouteProvider 路由提供者
func RouteProvider(r *gin.Engine) {
	Conf := bootstrap.GetAppConf()

	// 設置Middleware
	middleware.SetupMiddlewares(Conf.App.Site)

	// 全域Middleware載入
	r.Use(middleware.GlobalMiddlewares()...)

	// 載入Route
	switch Conf.App.Site {
	case "admin":
		// 專屬admin的route
	case "member":
		// 專屬member的route
	default:
		defaultR.LoadRoutes(r)
	}
}
