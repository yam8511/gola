package server

import (
	"sync"
)

// Run 啟動伺服器
func Run() {

	// 設置 router
	r := SetupRouter()

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
