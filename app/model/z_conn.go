package model

import (
	"errors"
	"gola/internal/database"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

// IModel Model有的func
type IModel interface {
	TableName() string
	Database(isMaster bool) databaseType
}

// 判斷資料不存在
func IsRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// 資料庫型態
type databaseType database.Type

// 建立資料庫連線
func (t databaseType) New() (*gorm.DB, error) {
	return database.GetPoolDB(database.Type(t), dbConf(t))
}

// 快取型態
type cacheType database.Type

// 建立快取連線
func (t cacheType) New() (*redis.Client, error) {
	return database.GetPoolRedis(database.Type(t), cacheConf(t))
}
