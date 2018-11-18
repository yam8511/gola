package database

import (
	"fmt"
	"gola/internal/bootstrap"

	"github.com/jinzhu/gorm"
)

// 取連線字串
func getConnectName(driver, host, port, database, username, password string) string {
	switch driver {
	case "mysql":
		return username + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?charset=utf8&parseTime=True"
	case "postgres":
		return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", host, port, database, username, password)
	}
	return ""
}

// NewOrmConnection 建立DB連線
func NewOrmConnection(master bool) (db *gorm.DB, err error) {
	Conf := bootstrap.GetAppConf()
	if master {
		connectName := getConnectName("mysql", Conf.DBMaster.Host, Conf.DBMaster.Port, Conf.DBMaster.DB, Conf.DBMaster.Username, Conf.DBMaster.Password)
		db, err = gorm.Open("mysql", connectName)
	} else {
		connectName := getConnectName("mysql", Conf.DBSlave.Host, Conf.DBSlave.Port, Conf.DBSlave.DB, Conf.DBSlave.Username, Conf.DBSlave.Password)
		db, err = gorm.Open("mysql", connectName)
	}
	if Conf.App.Env == "local" {
		db.LogMode(true)
	}
	return
}
