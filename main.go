package main

import (
	"gola/internal/entry"
	"os"
	// "gola/internal/database"
	// _ "github.com/jinzhu/gorm/dialects/mysql"
	// _ "github.com/jinzhu/gorm/dialects/postgres"
)

func init() {
	os.Setenv("TZ", "Asia/Taipei")
}

func main() {

	entry.Run(
		func() {
			// 設置資料表
			// model.SetupTable()
			// go database.SetupPool(150)
		},
	)
}
