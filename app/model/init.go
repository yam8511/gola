package model

import (
	"gola/internal/bootstrap"
	"gola/internal/database"
	"os"
)

// SetupTable 設置資料表
func SetupTable() {
	db, err := database.NewOrmConnection(true)
	if err != nil {
		bootstrap.WriteLog("ERROR", "DB連線失敗, "+err.Error())
		os.Exit(1)
		return
	}
	defer db.Close()

	modelList := []interface{}{
		new(User),
	}
	for i := range modelList {
		m := modelList[i]
		err = db.AutoMigrate(m).Error
		if err != nil {
			bootstrap.WriteLog("ERROR", "建立資料表失敗, "+err.Error())
			os.Exit(1)
			return
		}
	}

	env := bootstrap.GetAppConf().App.Env
	if env == "local" || env == "docker" {
		err = UserSeed(db)
		if err != nil {
			bootstrap.WriteLog("ERROR", "新增使用者資料失敗, "+err.Error())
			os.Exit(1)
			return
		}
	}
}
