package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	"gola/app/middleware"
	"gola/internal/bootstrap"
	"gola/internal/logger"

	"github.com/arl/statsviz"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// SetupRouter 配置路由
func SetupRouter() *gin.Engine {
	conf := bootstrap.GetAppConf()
	if !conf.App.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

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
	// healthz  健康檢測
	// config   預覽配置
	// metrics  撈取監控指標
	// pprof    分析效能
	r.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })
	r.GET("/config", func(c *gin.Context) { c.JSON(http.StatusOK, bootstrap.GetAppConf()) })
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/debug/pprof/*name", func(c *gin.Context) { pprofHandle[c.Param("name")].ServeHTTP(c.Writer, c.Request) })
	r.GET("/debug/statsviz/*name", func(c *gin.Context) {
		if strings.HasPrefix(c.Param("name"), "/ws") {
			statsviz.Ws(c.Writer, c.Request)
		} else {
			statsviz.Index.ServeHTTP(c.Writer, c.Request)
		}
	})

	return r
}

// CreateServer 建立伺服器
func CreateServer(router *gin.Engine) *http.Server {
	conf := bootstrap.GetAppConf().Server
	// 設定 Port
	var port string
	if conf.Port > 0 {
		port = ":" + strconv.Itoa(conf.Port)
	}
	addr := conf.IP + port

	// 建立 Server
	server := &http.Server{
		Addr:        addr,
		Handler:     router,
		ReadTimeout: 30 * time.Second,
		// WriteTimeout: 30 * time.Second,
		// MaxHeaderBytes: 1 << 20,
	}

	return server
}

// SignalListenAndServe 開啟Server & 系統信號監聽
func SignalListenAndServe(server *http.Server, waitFinish *sync.WaitGroup) {
	defer waitFinish.Done()
	defer func() {
		if err := recover(); err != nil {
			errMessage := fmt.Sprintf("Server 發生意外 Panic: %v", err)
			logger.Danger(errMessage)
			logger.Danger(string(debug.Stack()))
		}
	}()

	l, err := net.Listen("tcp", server.Addr)
	if err != nil {
		logger.Danger(fmt.Sprintf("Server 建立監聽連線失敗: %s", err.Error()))
		return
	}

	dl := NewDozListener(l, 0, true)

	// 嘗試使用 http2 server, 但是沒有效果, 只好使用一般https(內建http2功能)
	// err = http2.ConfigureServer(server, &http2.Server{})
	// if err != nil {
	// 	logger.Danger("轉成 http2 server 失敗: %s", err.Error())
	// 	return
	// }

	ctx, done := context.WithCancel(context.Background())
	go func() {
		var err error
		serverConf := bootstrap.GetAppConf().Server
		if serverConf.Secure {
			err = server.ServeTLS(dl, serverConf.TLS_Cert, serverConf.TLS_Key)
		} else {
			err = server.Serve(dl)
		}
		logger.Warn(fmt.Sprintf("🎃  Server 回傳 error (%v) 🎃", err))
		done()
	}()

	logger.Success("🐳  Web Server 開始服務! " + l.Addr().String() + "🐳")
	defer logger.Success("🔥  Web Server 結束服務!🔥")

	select {
	case <-ctx.Done():
	case <-bootstrap.GracefulDown():
		logger.Warn("⛔️  接受訊號 ⛔️")
	}

	select {
	case <-bootstrap.SingleFlightChan("Server.DozListener.Wait", func() (interface{}, error) {
		return nil, server.Shutdown(context.Background())
	}):
	case <-bootstrap.WaitOnceSignal():
		logger.Danger(`🚦  收到第二次訊號，強制結束 🚦`)
		os.Exit(2)
	}
}
