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
		ruleOptionData := map[string]interface{}{
			"event": "遊戲角色設定",
			"rule":  RuleOptions(),
		}
		連線.WriteJSON(ruleOptionData)

		rules := map[string]int{}
		for {
			_, msg, err := 連線.ReadMessage()
			if err != nil {
				return
			}

			err = json.Unmarshal(msg, &rules)
			if err != nil {
				連線.WriteJSON(ruleOptionData)
				return
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
		連線.WriteJSON(map[string]interface{}{
			"event": "遊戲建立成功",
			"token": 編號,
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
