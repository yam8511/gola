package model

import (
	errorcode "gola/app/common/errorcode"
	"gola/internal/bootstrap"
	"strings"
)

// SetupTable 設置資料表
func SetupTable() {

	modelList := []IModel{
		new(User),
	}

	env := bootstrap.GetAppConf().App.Env
	if env == "local" || env == "docker" {
		for _, m := range modelList {
			db, err := NewModelDB(m, true)
			if err != nil {
				errorcode.CheckGormConnError("DB連線失敗", err)
				return
			}

			err = db.AutoMigrate(m).Error
			if err != nil {
				errorcode.CheckGormConnError("建立資料表失敗", err)
				return
			}
		}
		return
	}

	missingTable := []string{}
	for _, m := range modelList {
		db, err := NewModelDB(m, false)
		if err != nil {
			errorcode.CheckGormConnError("DB連線失敗", err)
			return
		}

		if !db.HasTable(m.TableName()) {
			missingTable = append(missingTable, m.TableName())
		}
	}

	if len(missingTable) > 0 {
		errorcode.CheckGormConnError("缺少資料表: "+strings.Join(missingTable, ", "), nil)
		return
	}
}
