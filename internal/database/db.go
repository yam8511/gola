package database

import (
	"errors"
	"fmt"
	"gola/internal/bootstrap"
	"runtime"
	"time"

	"github.com/jinzhu/gorm"
)

// DB 資料庫
type DB struct {
	name            string            // 資料庫名稱
	config          *bootstrap.DBConf // 資料庫設定擋
	currentCount    int               // 用來記錄目前連線數
	connPool        []*gorm.DB        // 用來儲存連線
	getConnChannel  chan *gorm.DB     // 取連線
	putConnChannel  chan *gorm.DB     // 歸還連線
	askConnChannel  chan byte         // 要求連線
	closeChannel    chan byte         // 關閉通道
	getCountChannel chan int          // 取目前連線數量
	lastActionTime  time.Time         // 上次動作時間
}

// NewDB 建立DB
func NewDB(
	name string,
	config *bootstrap.DBConf,
) *DB {
	return &DB{
		name:   name,
		config: config,
	}
}

// PoolListen 連線池監聽
func (db *DB) PoolListen() {
	db.currentCount = 0
	db.connPool = []*gorm.DB{}
	db.getConnChannel = make(chan *gorm.DB)
	db.putConnChannel = make(chan *gorm.DB)
	db.askConnChannel = make(chan byte)
	db.closeChannel = make(chan byte)
	db.getCountChannel = make(chan int)
	db.lastActionTime = time.Now()
	ticker := time.NewTicker(time.Minute)
	isClosed := false

	for {
		select {
		case t := <-ticker.C:
			if t.Sub(db.lastActionTime) > time.Minute && !isClosed {
				count := len(db.connPool)
				bootstrap.WriteLog(
					"INFO",
					fmt.Sprintf(
						"⚡  [%s] 已經超過一分鐘沒有DB動作，所以關閉目前的連線，有%d條DB連線 ⚡",
						db.name,
						count,
					),
				)
				for i := 0; i < count; i++ {
					db.connPool[i].Close()
					db.connPool = []*gorm.DB{}
				}
			}
		case <-db.closeChannel:
			// 接收中斷訊號
			if !isClosed {
				isClosed = true
				bootstrap.WriteLog(
					"INFO",
					fmt.Sprintf(
						"🎃  [%s] DB接受訊號，準備關閉 🎃",
						db.name,
					),
				)
			}
		case conn := <-db.putConnChannel:
			// 歸還連線
			db.lastActionTime = time.Now()
			db.connPool = append(db.connPool, conn)
			db.currentCount--
		case <-db.askConnChannel:
			if db.currentCount < db.config.MaxConn && !isClosed {
				var conn *gorm.DB
				if len(db.connPool) > 0 {
					conn = db.connPool[0]
					db.connPool = db.connPool[1:]
				} else {
					var err error
					conn, err = NewOrmConnectionWithConf(db.config)
					if err != nil {
						bootstrap.WriteLog(
							"ERROR",
							fmt.Sprintf("[%s]  建立DB連線失敗! %s", db.name,
								err.Error(),
							),
						)
						continue
					}
				}
				db.currentCount++
				db.getConnChannel <- conn
			}
		case db.getCountChannel <- db.currentCount:
			// 有人來確認目前數量
		}
		runtime.Gosched()
		time.Sleep(time.Millisecond)
	}
}

// GetPoolConn 取連線
func (db *DB) GetPoolConn() (conn *gorm.DB, err error) {
	timeout := time.After(time.Second * 30)
	for {
		select {
		case <-timeout:
			err = errors.New("目前連線池數量不夠，請稍後再取DB連線")
			return
		case currentCount := <-db.getCountChannel:
			if currentCount < db.config.MaxConn {
				db.askConnChannel <- 0
				select {
				case conn = <-db.getConnChannel:
					return
				case <-time.After(time.Second * 30):
					err = errors.New("[" + db.name + "] 逾期30秒，請稍後再取DB連線")
					return
				}
			}
		}
		time.Sleep(time.Millisecond)
	}
}

// PutPoolConn  歸還連線
func (db *DB) PutPoolConn(conn *gorm.DB) {
	conn.Error = nil
	db.putConnChannel <- conn
}

// ClosePool 關閉連接池
func (db *DB) ClosePool() {
	defer func() {
		if catchErr := recover(); catchErr != nil {
			bootstrap.WriteLog(
				"WARNNING",
				fmt.Sprintf(
					"🎃  [%s] 關閉連接池發生意外, %v 🎃",
					db.name,
					catchErr,
				),
			)
		}
	}()
	close(db.closeChannel)
}

// WailPoolConnClosed 等待連線關閉
func (db *DB) WailPoolConnClosed() {
	defer bootstrap.WriteLog("INFO", "🔥 ["+db.name+"] 所有DB連線關閉結束  🔥")
	count := len(db.connPool)
	bootstrap.WriteLog("INFO", fmt.Sprintf(
		"⚡  [%s]等待所有DB連線關閉，目前有%d條連線 ⚡",
		db.name,
		count,
	))
	for {
		select {
		case count := <-db.getCountChannel:
			if count == 0 {
				return
			}
		}
		time.Sleep(time.Millisecond)
	}
}
