package database

import (
	"gola/internal/bootstrap"

	"github.com/jinzhu/gorm"
)

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

// NewOrmConnectionWithConf 用Conf建立新DB連線
func NewOrmConnectionWithConf(conf *bootstrap.DBConf) (db *gorm.DB, err error) {
	connectName := getConnectName(
		"mysql",
		conf.Host,
		conf.Port,
		conf.DB,
		conf.Username,
		conf.Password,
	)

	db, err = gorm.Open("mysql", connectName)
	if bootstrap.GetAppConf().App.Env == "local" {
		db.LogMode(true)
	}
	return
}
