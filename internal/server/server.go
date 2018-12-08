package server

import (
	"gola/app/model"
	"gola/internal/bootstrap"
	"gola/internal/database"
	"sync"
)

// Run 啟動伺服器
func Run() {
	// 載入設定檔資料
	Conf := bootstrap.LoadConfig()
	// 設置資料表
	model.SetupTable()
	go database.SetupPool(150)

	// 設置 router
	r := SetupRouter()

	// 建立 server
	server := CreateServer(r, Conf.App.Port, Conf.App.Host)
	waitFinish := new(sync.WaitGroup)

	// 系統信號監聽
	waitFinish.Add(1)
	go SignalListenAndServe(server, waitFinish)

	// 設置機器人監聽指令
	// waitFinish.Add(1)
	// go telegram.RunBot(waitFinish)

	// 等待結束
	waitFinish.Wait()
	database.WailDatabaseConnClosed()
}
