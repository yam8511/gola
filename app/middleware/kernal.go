package middleware

import (
	"gola/internal/bootstrap"
	"gola/internal/logger"
	"sync"

	"github.com/gin-gonic/gin"
)

var m = struct {
	sync.Once
	globalMiddlewares []gin.HandlerFunc          // 全域中介層群組
	middlewareGroups  []gin.HandlerFunc          // 中介層群組
	routeMiddleware   map[string]gin.HandlerFunc // 獨立的中介層
}{}

// 設置中介層
func setupMiddlewares() {
	// 全域中介層群組
	m.globalMiddlewares = []gin.HandlerFunc{
		gin.Recovery(),
	}
	if bootstrap.GetAppConf().App.Debug {
		m.globalMiddlewares = append(m.globalMiddlewares, gin.Logger())
	}

	// 中介層群組
	switch bootstrap.GetAppConf().App.Site {
	case "admin":
		m.middlewareGroups = []gin.HandlerFunc{}
	case "member":
		m.middlewareGroups = []gin.HandlerFunc{}
	default:
		m.middlewareGroups = []gin.HandlerFunc{}
	}

	// 獨立的中介層
	m.routeMiddleware = map[string]gin.HandlerFunc{
		"check_google_login": checkGoogleLogin,
	}
}

// GlobalMiddlewares 全域中介層群組
func GlobalMiddlewares() []gin.HandlerFunc {
	m.Do(setupMiddlewares)
	return m.globalMiddlewares
}

// GroupMiddlewares 取中介層群組
func GroupMiddlewares(name string) []gin.HandlerFunc {
	m.Do(setupMiddlewares)
	return m.middlewareGroups
}

// GetMiddleware 取獨立的中介層
func GetMiddleware(name string) gin.HandlerFunc {
	m.Do(setupMiddlewares)
	m, ok := m.routeMiddleware[name]
	if !ok || m == nil {
		m = func(c *gin.Context) {
			logger.Warn("Middleware doesn't exist [" + name + "]")
		}
	}
	return m
}
