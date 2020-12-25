package middleware

import (
	"gola/internal/bootstrap"
	"gola/internal/logger"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
		formatter := func(param gin.LogFormatterParams) string {
			if param.Latency > time.Minute {
				// Truncate in a golang < 1.8 safe way
				param.Latency = param.Latency - param.Latency%time.Second
			}

			if param.ErrorMessage == "" {
				param.ErrorMessage = "nil"
			}

			logrus.WithFields(logrus.Fields{
				// "req_time":  param.Keys["request_time"].(time.Time).Format(time.RFC3339Nano),
				// "res_time":  param.TimeStamp.Format(time.RFC3339Nano),
				"ip":        param.ClientIP,
				"execution": param.Latency,
				"status":    param.StatusCode,
				"path":      param.Path,
				"method":    param.Method,
				"error":     param.ErrorMessage,
			}).Infof("[GOLA-Request] %s", param.TimeStamp.Format(time.RFC3339Nano))
			return ""
		}

		m.globalMiddlewares = append(m.globalMiddlewares, gin.LoggerWithFormatter(formatter))
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
