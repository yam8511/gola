package model

import (
	"gola/internal/bootstrap"
	"gola/internal/database"

	"github.com/go-redis/redis"

	"github.com/jinzhu/gorm"
)

// 資料表名稱
const (
	/**
	 * ===============
	 *       DB
	 * ===============
	 */
	TableUser = "users"
)

// 資料庫常數。 ⭐️⭐️⭐️記得在下面的『 dbConf 』補上Config
const (
	DBMaster = database.Type("Default_master") // 一般資料庫(讀寫)
	DBSlave  = database.Type("Default_slave")  // 一般資料庫(唯讀)

	CacheMaster = database.Type("Cache_master") // 一般快取(讀寫)
	CacheSlave  = database.Type("Cache_slave")  // 一般快取(唯讀)
)

func dbConf(t database.Type) *bootstrap.DatabaseConf {
	switch t {
	case DBMaster:
		return bootstrap.GetAppConf().Databases.DefaultMaster
	case DBSlave:
		return bootstrap.GetAppConf().Databases.DefaultSlave
	default:
		panic("DB型態錯誤")
	}
}

// cacheConf DB的設定檔
func cacheConf(t database.Type) *bootstrap.CacheConf {
	switch t {
	case CacheMaster:
		return bootstrap.GetAppConf().Caches.DefaultMaster
	case CacheSlave:
		return bootstrap.GetAppConf().Caches.DefaultSlave
	default:
		panic("Cache型態錯誤")
	}
}

// IModel Model有的func
type IModel interface {
	TableName() string
	Database(isMaster bool) database.Type
}

// NewModelDB 建立新的Model的DB連線
func NewModelDB(m IModel, master bool) (*gorm.DB, error) {
	dbType := m.Database(master)
	return NewDB(dbType)
}

// NewDB 建立新的DB連線
func NewDB(dbType database.Type) (*gorm.DB, error) {
	return database.GetPoolDB(dbType, dbConf(dbType))
}

// NewRedis 建立新的Redis連線
func NewRedis(cacheType database.Type) (*redis.Client, error) {
	return database.GetPoolRedis(cacheType, cacheConf(cacheType))
}
