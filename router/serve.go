package router

import (
	"fmt"
	"gola/internal/bootstrap"
	"gola/router/provider"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置路由
func SetupRouter() (r *gin.Engine) {
	if bootstrap.LoadConfig().App.Debug {
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
		r.Use(gin.Recovery())
	}

	provider.RouteProvider(r)

	return r
}

// CreateServer 建立伺服器
func CreateServer(router *gin.Engine, port, host string, args ...string) *http.Server {
	// 建立 Server
	server := &http.Server{
		Addr:    port,
		Handler: router,
		// ReadTimeout:  5 * time.Second,
		// WriteTimeout: 10 * time.Second,
		// MaxHeaderBytes: 1 << 20,
	}

	return server
}

// SignalListenAndServe 開啟Server & 系統信號監聽
func SignalListenAndServe(server *http.Server, waitFinish *sync.WaitGroup) {
	defer waitFinish.Done()
	defer func() {
		if err := recover(); err != nil {
			errMessage := fmt.Sprintf("❌  Server 發生意外 Error: %v ❌", err)
			bootstrap.WriteLog("ERROR", errMessage)
		}
	}()

	wg := make(chan int, 2)

	go func() {
		err := server.ListenAndServe()
		bootstrap.WriteLog("WARNING", fmt.Sprintf("🎃  Server 回傳 error (%v) 🎃", err))
		wg <- 1
	}()

	go func() {
		receivedSignal := <-bootstrap.GracefulDown()
		bootstrap.WriteLog("INFO", fmt.Sprintf("🎃  接受訊號 <- %v 🎃", receivedSignal))
		wg <- 0
	}()

	bootstrap.WriteLog("INFO", "🐳  Web Server 開始服務! 🐳")
	defer bootstrap.WriteLog("INFO", "🔥  Web Server 結束服務!🔥")
	select {
	case <-wg:
	}
}
