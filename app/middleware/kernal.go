package middleware

import (
	"gola/internal/bootstrap"

	"github.com/gin-gonic/gin"
)

// 全域中介層群組
var globalMiddlewares = []gin.HandlerFunc{
	gin.Recovery(),
}

// 中介層群組
var middlewareGroups = []gin.HandlerFunc{}

// 獨立的中介層
var routeMiddleware = map[string]gin.HandlerFunc{}

// SetupMiddlewares 設置中介層
func SetupMiddlewares(site string) {
	// 中介層群組
	switch site {
	case "admin":
		middlewareGroups = []gin.HandlerFunc{}
	case "member":
		middlewareGroups = []gin.HandlerFunc{}
	default:
		middlewareGroups = []gin.HandlerFunc{}
	}

	// 獨立的中介層
	routeMiddleware = map[string]gin.HandlerFunc{
		"check_google_login": checkGoogleLogin,
	}
}

// GlobalMiddlewares 全域中介層群組
func GlobalMiddlewares() []gin.HandlerFunc {
	if globalMiddlewares == nil {
		globalMiddlewares = []gin.HandlerFunc{}
	}
	return globalMiddlewares
}

// GroupMiddlewares 取中介層群組
func GroupMiddlewares(name string) []gin.HandlerFunc {
	return middlewareGroups
}

// GetMiddleware 取獨立的中介層
func GetMiddleware(name string) gin.HandlerFunc {
	m, ok := routeMiddleware[name]
	if !ok || m == nil {
		m = func(c *gin.Context) {
			bootstrap.WriteLog("WARNING", "Middleware doesn't exist ["+name+"]")
		}
	}
	return m
}
