package database

import (
	"gola/internal/bootstrap"
	"strconv"
	"time"

	redis "github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewOrmConnectionWithConf 用Conf建立新DB連線
func NewOrmConnectionWithConf(conf *bootstrap.DatabaseConf) (db *gorm.DB, err error) {
	dsn := getConnectName(
		"mysql",
		conf.Host,
		conf.Port,
		conf.DB,
		conf.Username,
		conf.Password,
	)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if bootstrap.GetAppConf().App.Env == "local" {
		db = db.Debug()
	}
	return
}

// NewRedisConnWithConf 用設定建立Redis連線
func NewRedisConnWithConf(conf *bootstrap.CacheConf) *redis.Client {
	addr := conf.Host
	if conf.Port > 0 {
		addr += ":" + strconv.Itoa(conf.Port)
	}
	redisConn := redis.NewClient(&redis.Options{
		Addr:        addr,
		Password:    conf.Password,
		DB:          conf.DB, // use default DB
		PoolSize:    conf.MaxConn,
		IdleTimeout: time.Minute,
	})
	return redisConn
}
