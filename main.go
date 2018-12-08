package main

import (
	"gola/internal/entry"
	"os"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	// _ "github.com/jinzhu/gorm/dialects/postgres"
)

func init() {
	os.Setenv("TZ", "Asia/Taipei")
}

func main() {
	entry.Run()
}
