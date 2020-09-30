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
				errorcode.Code_Undefined.New("DB連線失敗: %w", err)
				return
			}

			err = db.AutoMigrate(m).Error
			if err != nil {
				errorcode.Code_Undefined.New("DB連線失敗: %w", err)
				return
			}
		}
		return
	}

	missingTable := []string{}
	for _, m := range modelList {
		db, err := NewModelDB(m, false)
		if err != nil {
			errorcode.Code_Undefined.New("DB連線失敗: %w", err)
			return
		}

		if !db.HasTable(m.TableName()) {
			missingTable = append(missingTable, m.TableName())
		}
	}

	if len(missingTable) > 0 {
		errorcode.Code_Param_Required.New("缺少資料表: " + strings.Join(missingTable, ", "))
		return
	}
}
