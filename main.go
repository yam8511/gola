package main

import (
	"fmt"
	"gola/app/model"
	"gola/internal/bootstrap"
	"gola/router"
	"os"
	"sync"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	// _ "github.com/jinzhu/gorm/dialects/postgres"
)

func init() {
	os.Setenv("TZ", "Asia/Taipei")
	// 顯示說明
	Usage()
}

// Usage 顯示說明
func Usage() {
	if bootstrap.GetAppEnv() == "" {
		fmt.Printf(`
			📖 APP 說明 📖
			需傳入以下環境變數：

			⚙  APP_ENV : 專案環境
				✏ docker 容器開發
				✏ local 本機開發
				✏ qatest 測試
				✏ prod 正式

			📌  舉例： APP_ENV=local ./server
`)
		os.Exit(0)
	} else {
		bootstrap.WriteLog("INFO", fmt.Sprintf("⚙  APP_ENV: %s", bootstrap.GetAppEnv()))
	}
}

// @title API文件
// @version 0.1.0
// @description  API文件範例
// @termsOfService http://swagger.io/terms/
// @license.name Zuolar
func main() {
	// 設定優雅結束程序
	bootstrap.SetupGracefulSignal()
	// 載入設定檔資料
	Conf := bootstrap.LoadConfig()
	// 設置資料表
	model.SetupTable()

	// 設置 router
	r := router.SetupRouter()

	// 建立 server
	server := router.CreateServer(r, Conf.App.Port, Conf.App.Host)
	waitFinish := new(sync.WaitGroup)

	// 系統信號監聽
	waitFinish.Add(1)
	go router.SignalListenAndServe(server, waitFinish)

	// 設置機器人監聽指令
	// waitFinish.Add(1)
	// go telegram.RunBot(waitFinish)

	// 等待結束
	waitFinish.Wait()
}
