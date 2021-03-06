package dance

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// EnterGame 進入遊戲
func EnterGame(連線 *websocket.Conn, 編號 string) {
	var 遊戲 *Game
	var 新進入 bool
	if 編號 == "" {
		編號, 遊戲 = 新遊戲()
		新進入 = true
	} else {
		遊戲 = 取遊戲(編號)
	}

	// 檢查遊戲是否設定了
	if !遊戲.HasSetup() {
		setupOptionData := TransferData{
			Sound:  "請輸入卡牌設定",
			Action: 牌組設定,
			Data: map[string]interface{}{
				"BasicCard": 基本牌組,
				"Combine":   規則牌組設定,
			},
		}

		err := 連線.WriteJSON(setupOptionData)
		if err != nil {
			return
		}

		setupConf := 規則設定{}
		for {
			so, err := waitSocketBack(連線, 牌組設定)
			if err != nil {
				return
			}

			if 遊戲.HasSetup() {
				break
			}

			var ok bool
			setupConf, ok = 規則牌組設定[string(so.Reply)]
			if !ok {
				err = json.Unmarshal([]byte(so.Reply), &setupConf)
				if err != nil {
					err = 連線.WriteJSON(setupOptionData)
					if err != nil {
						return
					}
					continue
				}
			}
			break
		}

		if 遊戲.HasSetup() {
			新進入 = false
		} else {
			遊戲.Setup(setupConf)
		}
	}

	if 新進入 {
		err := 連線.WriteJSON(TransferData{
			Sound:   "可分享序號一起玩",
			Display: "可分享序號一起玩: " + 編號,
			Action:  拿到Token,
			Data:    編號,
		})

		if err != nil {
			return
		}

		for {
			_, err := waitSocketBack(連線, 拿到Token)
			if err != nil {
				return
			}

			break
		}

	}

	遊戲.加入(連線)
}

var 遊戲管理中心 = struct {
	sync.Once
	遊戲群 map[string]*Game
	讀寫鎖 *sync.RWMutex
}{}

func 建立新遊戲(編號 string) *Game {
	遊戲管理中心.Do(func() {
		遊戲管理中心.讀寫鎖 = &sync.RWMutex{}
		遊戲管理中心.遊戲群 = map[string]*Game{}
	})

	game := &Game{}
	遊戲管理中心.讀寫鎖.Lock()
	遊戲管理中心.遊戲群[編號] = game
	遊戲管理中心.讀寫鎖.Unlock()

	return game
}

// 新遊戲 建立新的一場遊戲
func 新遊戲() (string, *Game) {
	編號 := fmt.Sprintf("%x", md5.Sum([]byte(uuid.New().String())))
	遊戲 := 建立新遊戲(編號)
	return 編號, 遊戲
}

// 取遊戲 建立新的一場遊戲
func 取遊戲(編號 string) *Game {
	遊戲管理中心.Do(func() {
		遊戲管理中心.讀寫鎖 = &sync.RWMutex{}
		遊戲管理中心.遊戲群 = map[string]*Game{}
	})
	遊戲管理中心.讀寫鎖.RLock()
	遊戲, 存在 := 遊戲管理中心.遊戲群[編號]
	遊戲管理中心.讀寫鎖.RUnlock()

	遊戲不存在 := !存在
	if 遊戲不存在 {
		遊戲 = 建立新遊戲(編號)
	}

	return 遊戲
}
