package server

import (
	"sync"

	"github.com/gin-gonic/gin"
)

// Run 啟動伺服器
func Run(provider func(*gin.Engine)) {

	// 設置 router
	r := SetupRouter()

	// 自定義route
	if provider != nil {
		provider(r)
	}

	// 建立 server
	server := CreateServer(r)
	waitFinish := new(sync.WaitGroup)

	// 系統信號監聽
	waitFinish.Add(1)
	go SignalListenAndServe(server, waitFinish)

	// 設置機器人監聽指令
	// waitFinish.Add(1)
	// go telegram.RunBot(waitFinish)

	// 等待結束
	waitFinish.Wait()
	// database.WailDatabaseConnClosed()
}
