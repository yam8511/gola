package main

import (
	"log"
	"os"
	"strconv"

	"github.com/gorilla/websocket"
	// "gola/internal/entry"
	// "gola/internal/database"
	// _ "github.com/jinzhu/gorm/dialects/mysql"
	// _ "github.com/jinzhu/gorm/dialects/postgres"
)

func init() {
	os.Setenv("TZ", "Asia/Taipei")
}

func main() {

	// entry.Run(
	// 	func() {
	// 		// 設置資料表
	// 		// model.SetupTable()
	// 		// go database.SetupPool(150)
	// 	},
	// )

	// http://www.87g.com/lrs/59119.html
	玩家們 := []玩家{}
	for i := 0; i < 3; i++ {
		玩家們 = append(玩家們, 新的玩家(i+1, nil))
		log.Printf("%+v", 玩家們[i])
	}
	log.Println(玩家們)
}

func 新的玩家(位子 int, 連線 *websocket.Conn) *平民 {
	return &平民{
		位子: 位子,
		連線: 連線,
	}
}

type 平民 struct {
	位子    int
	開眼睛中  bool
	被投票狀態 bool
	出局狀態  bool
	連線    *websocket.Conn
}

func (我 *平民) 號碼() int {
	return 我.位子
}

func (我 *平民) 閉眼睛() {
	我.開眼睛中 = false
}

func (我 *平民) 開眼睛() {
	我.開眼睛中 = true
}

func (我 *平民) 投票() (int, error) {
	if 我.連線 == nil {
		return 0, nil
	}

	_, msg, err := 我.連線.ReadMessage()
	if err != nil {
		return 0, nil
	}

	no := 0
	no, err = strconv.Atoi(string(msg))
	return no, err
}

func (我 *平民) 被投票(是 bool) {
	我.被投票狀態 = 是
}

func (我 *平民) 被投票了() bool {
	return 我.被投票狀態
}

func (我 *平民) 出局了() bool {
	return 我.出局狀態
}

func (我 *平民) 種族() 種族 {
	return 人質
}

func (我 *平民) 換號碼(新位子 int) int {
	舊的位子 := 我.位子
	我.位子 = 新位子
	return 舊的位子
}
