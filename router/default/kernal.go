package router

import "github.com/gin-gonic/gin"

// LoadRoutes 載入 routes
func LoadRoutes(r *gin.Engine) {
	LoadWebRouter(r)
	LoadAPIRouter(r)
	LoadGraphQLRouter(r)
}
