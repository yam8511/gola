package model

import (
	"fmt"
	"gola/app/common/errorcode"
	"gola/internal/bootstrap"
	"strings"
)

// SetupTable 設置資料表
func SetupTable(modelList ...IModel) {

	env := bootstrap.GetAppConf().App.Env
	if env == "local" {
		for _, m := range modelList {
			db, err := m.Database(true).New()
			if err != nil {
				panic(errorcode.CheckGormConnError(errorcode.Code_Undefined, fmt.Errorf("DB連線失敗: %w", err)))
			}
			err = db.AutoMigrate(m)
			if err != nil {
				panic(errorcode.CheckGormConnError(errorcode.Code_Undefined, fmt.Errorf("建立資料表失敗: %w", err)))
			}
		}
		return
	}

	missingTable := []string{}
	for _, m := range modelList {
		db, err := m.Database(true).New()
		if err != nil {
			panic(errorcode.CheckGormConnError(errorcode.Code_Undefined, fmt.Errorf("DB連線失敗: %w", err)))
		}
		db = db.Debug()
		if !db.Migrator().HasTable(m.TableName()) {
			missingTable = append(missingTable, m.TableName())
		}
	}

	if len(missingTable) > 0 {
		panic(errorcode.CheckGormConnError(errorcode.Code_Undefined, fmt.Errorf("缺少資料表: "+strings.Join(missingTable, ", "))))
	}
}
