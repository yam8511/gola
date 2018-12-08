package middleware

import (
	"gola/internal/bootstrap"

	"github.com/gin-gonic/gin"
)

// 中介層群組
var middlewareGroups = map[string][]gin.HandlerFunc{}

// 獨立的中介層
var routeMiddleware = map[string]gin.HandlerFunc{
	"check_google_login": checkGoogleLogin,
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

// GetMiddlewareGroup 取中介層群組
func GetMiddlewareGroup(name string) []gin.HandlerFunc {
	m, ok := middlewareGroups[name]
	if !ok || m == nil {
		bootstrap.WriteLog("WARNING", "Middleware Group doesn't exist ["+name+"]")
		m = []gin.HandlerFunc{}
	}
	return m
}
