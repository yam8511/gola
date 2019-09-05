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
var middlewareGroups = map[string][]gin.HandlerFunc{}

// 獨立的中介層
var routeMiddleware = map[string]gin.HandlerFunc{
	"check_google_login": checkGoogleLogin,
}

// GlobalMiddlewares 全域中介層群組
func GlobalMiddlewares() []gin.HandlerFunc {
	if globalMiddlewares == nil {
		globalMiddlewares = []gin.HandlerFunc{}
	}
	return globalMiddlewares
}

// MiddlewareGroup 取中介層群組
func MiddlewareGroup(name string) []gin.HandlerFunc {
	m, ok := middlewareGroups[name]
	if !ok || m == nil {
		bootstrap.WriteLog("WARNING", "Middleware Group doesn't exist ["+name+"]")
		m = []gin.HandlerFunc{}
	}
	return m
}

// Middleware 取獨立的中介層
func Middleware(name string) gin.HandlerFunc {
	m, ok := routeMiddleware[name]
	if !ok || m == nil {
		m = func(c *gin.Context) {
			bootstrap.WriteLog("WARNING", "Middleware doesn't exist ["+name+"]")
		}
	}
	return m
}
