package server

import (
	"context"
	"fmt"
	"gola/internal/bootstrap"
	"gola/internal/logger"
	"gola/router"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置路由
func SetupRouter() (r *gin.Engine) {
	if !bootstrap.GetAppConf().App.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r = gin.New()
	router.RouteProvider(r)

	return r
}

// CreateServer 建立伺服器
func CreateServer(router *gin.Engine) *http.Server {
	conf := bootstrap.GetAppConf().Server
	// 設定 Port
	var port = conf.Port
	if conf.AutoPort && os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	if port != "" {
		port = ":" + port
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

	go func() {
		err := server.Serve(l)
		logger.Warn(fmt.Sprintf("🎃  Server 回傳 error (%v) 🎃", err))
	}()

	logger.Success("🐳  Web Server 開始服務! " + l.Addr().String() + "🐳")
	defer logger.Success("🔥  Web Server 結束服務!🔥")

	<-bootstrap.GracefulDown()
	logger.Warn("🎃  接受訊號 🎃")

	select {
	case <-bootstrap.SingleFlightChan("Server.DozListener.Wait", func() (interface{}, error) {
		return nil, server.Shutdown(context.Background())
	}):
	case <-bootstrap.WaitOnceSignal():
		logger.Danger(`🚦  收到第二次訊號，強制結束 🚦`)
		os.Exit(2)
	}
}
