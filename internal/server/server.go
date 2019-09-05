package server

import (
	"gola/internal/bootstrap"
	"os"
	"sync"
)

// Run 啟動伺服器
func Run() {
	// 載入設定檔資料
	Conf := bootstrap.GetAppConf()

	// 設置 router
	r := SetupRouter()

	// 設定 Port
	var port = Conf.App.Port
	if Conf.App.AutoPort && os.Getenv("PORT") != "" {
		port = ":" + os.Getenv("PORT")
	}

	// 建立 server
	server := CreateServer(r, port, Conf.App.Host)
	waitFinish := new(sync.WaitGroup)

	// 系統信號監聽
	waitFinish.Add(1)
	go SignalListenAndServe(server, waitFinish, Conf.App.MaxConn)

	// 設置機器人監聽指令
	// waitFinish.Add(1)
	// go telegram.RunBot(waitFinish)

	// 等待結束
	waitFinish.Wait()
	// database.WailDatabaseConnClosed()
}
