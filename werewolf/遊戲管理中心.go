package werewolf

import (
	"encoding/json"
	"log"
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
	if !遊戲.初始設定過() {
		ruleOptionData := 傳輸資料{
			Sound:  "請輸入角色設定",
			Action: 角色設定,
			Data: map[string]interface{}{
				"rule":    RuleOptions(),
				"combine": []string{"4", "5", "6"},
			},
		}
		連線.WriteJSON(ruleOptionData)

		rules := map[RULE]int{}
		for {
			_, msg, err := 連線.ReadMessage()
			if err != nil {
				return
			}

			switch string(msg) {
			case "4":
				rules = map[RULE]int{
					平民: 2,
					狼人: 1,
					騎士: 1,
				}
			case "5":
				rules = map[RULE]int{
					平民: 3,
					狼人: 1,
					騎士: 1,
				}
			case "6":
				rules = map[RULE]int{
					平民: 2,
					狼人: 2,
					騎士: 2,
				}
			default:
				err = json.Unmarshal(msg, &rules)
				if err != nil {
					連線.WriteJSON(ruleOptionData)
					continue
				}
			}

			break
		}

		if 遊戲.初始設定過() {
			新進入 = false
		} else {
			遊戲.初始設定(rules)
		}
	}

	if 新進入 {
		連線.WriteJSON(傳輸資料{
			Sound:  "可分享序號一起玩",
			Action: 新序號,
			Data:   編號,
		})

		for {
			_, _, err := 連線.ReadMessage()
			if err != nil {
				log.Print("等待開始錯誤 => ", err)
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
	// http://www.87g.com/lrs/59119.html
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
	編號 := uuid.New().String()
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
