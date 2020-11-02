package router

import (
	"gola/app/middleware"
	"gola/internal/bootstrap"
	"net/http"
	"net/http/pprof"
	"strings"

	"github.com/arl/statsviz"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// RouteProvider 路由提供者
func RouteProvider(r *gin.Engine) {
	pprofHandle := map[string]http.HandlerFunc{
		"/allocs":       pprof.Handler("allocs").ServeHTTP,
		"/block":        pprof.Handler("block").ServeHTTP,
		"/goroutine":    pprof.Handler("goroutine").ServeHTTP,
		"/heap":         pprof.Handler("heap").ServeHTTP,
		"/mutex":        pprof.Handler("mutex").ServeHTTP,
		"/threadcreate": pprof.Handler("threadcreate").ServeHTTP,
		"/cmdline":      pprof.Cmdline,
		"/profile":      pprof.Profile,
		"/symbol":       pprof.Symbol,
		"/trace":        pprof.Trace,
		"/":             pprof.Index,
	}

	// 全域Middleware載入
	r.Use(middleware.GlobalMiddlewares()...)

	// 設置全域Route
	// healthz    健康檢測
	// config     預覽配置
	// metrics    撈取監控指標
	// pprof      分析效能
	// statsviz   即時效能
	r.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })
	r.GET("/config", func(c *gin.Context) { c.JSON(http.StatusOK, bootstrap.GetAppConf()) })
	r.Any("/metrics", gin.WrapH(promhttp.Handler()))
	r.Any("/debug/pprof/*name", func(c *gin.Context) { pprofHandle[c.Param("name")].ServeHTTP(c.Writer, c.Request) })
	r.GET("/debug/statsviz/*name", func(c *gin.Context) {
		if strings.HasPrefix(c.Param("name"), "/ws") {
			statsviz.Ws(c.Writer, c.Request)
		} else {
			statsviz.Index.ServeHTTP(c.Writer, c.Request)
		}
	})
}
