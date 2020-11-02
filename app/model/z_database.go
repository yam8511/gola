package model

import (
	"gola/internal/bootstrap"
)

// 資料庫常數
const (
	DB_Master = databaseType("Default_master") // 一般資料庫(讀寫)
	DB_Slave  = databaseType("Default_slave")  // 一般資料庫(唯讀)
)

// 資料庫設定檔
func dbConf(t databaseType) *bootstrap.DatabaseConf {
	switch t {
	case DB_Master:
		return bootstrap.GetAppConf().Databases.DefaultMaster
	case DB_Slave:
		return bootstrap.GetAppConf().Databases.DefaultSlave
	default:
		panic("DB型態錯誤")
	}
}

// 快取常數
const (
	Cache_Master = cacheType("Default_master") // 一般資料庫(讀寫)
	Cache_Slave  = cacheType("Default_slave")  // 一般資料庫(唯讀)
)

// 快取設定檔
func cacheConf(t cacheType) *bootstrap.CacheConf {
	switch t {
	case Cache_Master:
		return bootstrap.GetAppConf().Caches.DefaultMaster
	case Cache_Slave:
		return bootstrap.GetAppConf().Caches.DefaultSlave
	default:
		panic("Cache型態錯誤")
	}
}
