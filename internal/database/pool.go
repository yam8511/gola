package database

import (
	"errors"
	"fmt"
	"gola/internal/bootstrap"
	"runtime"
	"time"

	"github.com/jinzhu/gorm"
)

var (
	masterPool        []*gorm.DB
	slavePool         []*gorm.DB
	masterGetChannel  chan *gorm.DB
	masterPutChannel  chan *gorm.DB
	slaveGetChannel   chan *gorm.DB
	slavePutChannel   chan *gorm.DB
	masterAskChannel  chan byte
	slaveAskChannel   chan byte
	dbCountChannel    chan int
	currenctConnCount int
	maxConnCount      int
)

// SetupPool 設置連接池
func SetupPool(max int) {
	maxConnCount = max
	currenctConnCount = 0
	masterPool = []*gorm.DB{}
	slavePool = []*gorm.DB{}
	masterGetChannel = make(chan *gorm.DB)
	masterPutChannel = make(chan *gorm.DB)
	slaveGetChannel = make(chan *gorm.DB)
	slavePutChannel = make(chan *gorm.DB)
	masterAskChannel = make(chan byte)
	slaveAskChannel = make(chan byte)
	dbCountChannel = make(chan int)
	sig := bootstrap.WaitOnceSignal()
	catchSignal := false
	lastActionTime := time.Now()
	ticker := time.NewTicker(time.Minute)

	for {
		select {
		case t := <-ticker.C:
			if t.Sub(lastActionTime) > time.Minute && !catchSignal {
				masterCount := len(masterPool)
				slaveCount := len(slavePool)
				bootstrap.WriteLog("INFO", fmt.Sprintf("⚡  已經超過一分鐘DB沒有動作，所以關閉目前的連線，有%d條master連線, %d條slave連線 ⚡", masterCount, slaveCount))
				for i := 0; i < masterCount; i++ {
					masterPool[i].Close()
					masterPool = []*gorm.DB{}
				}
				for i := 0; i < slaveCount; i++ {
					slavePool[i].Close()
					slavePool = []*gorm.DB{}
				}
			}
		case s := <-sig:
			// 接收中斷訊號
			if !catchSignal {
				bootstrap.WriteLog("INFO", fmt.Sprintf("🎃  DB接受訊號 <- %v 🎃", s))
				catchSignal = true
			}
		case db := <-masterPutChannel:
			// 歸還Master連線
			lastActionTime = time.Now()
			masterPool = append(masterPool, db)
			currenctConnCount--
		case db := <-slavePutChannel:
			// 歸還Slave連線
			lastActionTime = time.Now()
			slavePool = append(slavePool, db)
			currenctConnCount--
		case <-masterAskChannel:
			// 要求Master連線
			if currenctConnCount < maxConnCount && !catchSignal {
				var db *gorm.DB
				if len(masterPool) > 0 {
					db = masterPool[0]
					masterPool = masterPool[1:]
				} else {
					var err error
					db, err = NewOrmConnection(true)
					if err != nil {
						bootstrap.WriteLog("ERROR", "建立DB連線失敗! "+err.Error())
						continue
					}
				}
				currenctConnCount++
				masterGetChannel <- db
			}
		case <-slaveAskChannel:
			// 要求Slave連線
			if currenctConnCount < maxConnCount && !catchSignal {
				var db *gorm.DB
				if len(slavePool) > 0 {
					db = slavePool[0]
					slavePool = slavePool[1:]
				} else {
					var err error
					db, err = NewOrmConnection(false)
					if err != nil {
						bootstrap.WriteLog("ERROR", "建立DB連線失敗! "+err.Error())
						continue
					}
				}
				currenctConnCount++
				slaveGetChannel <- db
			}
		case dbCountChannel <- currenctConnCount:
			// 有人來確認目前數量
		}
		runtime.Gosched()
		time.Sleep(time.Millisecond)
	}
}

// WailDatabaseConnClosed 等待連線關閉
func WailDatabaseConnClosed() {
	defer bootstrap.WriteLog("INFO", "🔥  所有DB連線關閉結束  🔥")
	masterCount := len(masterPool)
	slaveCount := len(slavePool)
	bootstrap.WriteLog("INFO", fmt.Sprintf("⚡  等待所有DB連線關閉，目前有%d條master連線, %d條slave連線 ⚡", masterCount, slaveCount))
	for {
		select {
		case count := <-dbCountChannel:
			if count == 0 {
				return
			}
		}
		time.Sleep(time.Millisecond)
	}
}

// GetPoolMasterDB  取master連線池的DB
func GetPoolMasterDB() (db *gorm.DB, err error) {
	timeout := time.After(time.Second * 30)
	for {
		select {
		case <-timeout:
			err = errors.New("目前連線池數量不夠，請稍後再取DB連線")
			return
		case currentCount := <-dbCountChannel:
			if currentCount < maxConnCount {
				masterAskChannel <- 0
				select {
				case db = <-masterGetChannel:
					return
				case <-time.After(time.Second * 30):
					err = errors.New("逾期30秒，請稍後再取DB連線")
					return
				}
			}
		}
		time.Sleep(time.Millisecond)
	}
}

// GetPoolSlaveDB  取slave連線池的DB
func GetPoolSlaveDB() (db *gorm.DB, err error) {
	timeout := time.After(time.Second * 30)
	for {
		select {
		case <-timeout:
			err = errors.New("目前連線池數量不夠，請稍後再取DB連線")
			return
		case currentCount := <-dbCountChannel:
			if currentCount < maxConnCount {
				slaveAskChannel <- 0
				select {
				case db = <-slaveGetChannel:
					return
				case <-time.After(time.Second * 30):
					err = errors.New("逾期30秒，請稍後再取DB連線")
					return
				}
			}
		}
		time.Sleep(time.Millisecond)
	}
}

// PutPoolMasterDB  歸還master連線
func PutPoolMasterDB(db *gorm.DB) {
	putPoolDB(db, true)
}

// PutPoolSlaveDB  歸還slave連線
func PutPoolSlaveDB(db *gorm.DB) {
	putPoolDB(db, false)
}

func putPoolDB(db *gorm.DB, master bool) {
	db.Error = nil
	if master {
		masterPutChannel <- db
	} else {
		slavePutChannel <- db
	}
}
