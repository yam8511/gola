package database

import (
	"gola/internal/bootstrap"
	"sync"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
)

// Type 資料庫型態
type Type string

// 預設資料庫型態
const (
	DefaultDatabaseType = Type("")
	DefaultCacheType    = Type("")
)

var poolDB map[Type]*gorm.DB
var mxDB *sync.RWMutex
var poolCache map[Type]*redis.Client
var mxCache *sync.RWMutex

func init() {
	poolDB = map[Type]*gorm.DB{}
	mxDB = &sync.RWMutex{}
	poolCache = map[Type]*redis.Client{}
	mxCache = &sync.RWMutex{}
}

// GetPoolDB 取連線池的DB
func GetPoolDB(t Type, conf *bootstrap.DatabaseConf) (*gorm.DB, error) {
	mxDB.RLock()
	db, ok := poolDB[t]
	mxDB.RUnlock()
	if ok {
		return db, nil
	}

	db, err := NewOrmConnectionWithConf(conf)
	if err != nil {
		return nil, err
	}
	mxDB.Lock()
	poolDB[t] = db
	mxDB.Unlock()
	return db, nil
}

// GetPoolRedis 取連線池的DB
func GetPoolRedis(t Type, conf *bootstrap.CacheConf) (*redis.Client, error) {
	mxCache.RLock()
	db, ok := poolCache[t]
	mxCache.RUnlock()
	if ok {
		return db, nil
	}

	db = NewRedisConnWithConf(conf)
	mxCache.Lock()
	poolCache[t] = db
	mxCache.Unlock()
	return db, nil
}
