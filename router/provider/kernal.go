package provider

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// RouteProvider 路由提供者
func RouteProvider(r *gin.Engine) {
	LoadWebRouter(r)
	LoadAPIRouter(r)
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
