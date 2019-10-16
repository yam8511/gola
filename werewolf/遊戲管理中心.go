package werewolf

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"sync"

	datastruct "gola/app/common/data_struct"

	"github.com/google/uuid"
)

// EnterGame 進入遊戲
func EnterGame(連線 *datastruct.WebSocketConn, 編號 string) {
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
		combine := 快速組合()
		ruleOptionData := 傳輸資料{
			Sound:  "請輸入角色設定",
			Action: 角色設定,
			Data: map[string]interface{}{
				"rule":    角色選單(),
				"combine": combine,
			},
		}

		closed := 連線.Write(ruleOptionData)
		if closed {
			return
		}

		rules := map[RULE]int{}
		for {
			so, err := waitSocketBack(連線.Conn, 角色設定)
			if err != nil {
				return
			}

			if 遊戲.初始設定過() {
				break
			}

			var ok bool
			rules, ok = combine[string(so.Reply)]
			if !ok {
				err = json.Unmarshal([]byte(so.Reply), &rules)
				if err != nil {
					closed = 連線.Write(ruleOptionData)
					if closed {
						return
					}
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
		closed := 連線.Write(傳輸資料{
			Sound:   "可分享序號一起玩",
			Display: "可分享序號一起玩: " + 編號,
			Action:  拿到Token,
			Data:    編號,
		})

		if closed {
			return
		}

		for {
			_, err := waitSocketBack(連線.Conn, 拿到Token)
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
