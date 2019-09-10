package werewolf

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Game sera找明俊玩遊戲
type Game struct {
	房主  *websocket.Conn
	讀寫鎖 sync.RWMutex	玩家們 []Player
	連線池 []*websocket.Conn
	階段  階段
	通訊  chan int
}

// Join 加入遊戲
func (遊戲 *Game) Join(連線 *websocket.Conn) {
	遊戲.讀寫鎖.RLock()
	遊戲已經開始 := 遊戲.階段 == 開始階段
	遊戲.讀寫鎖.RUnlock()
	if 遊戲已經開始 {
		連線.WriteJSON("遊戲已經開始")
		return
	}

	遊戲.儲存連線(連線)

	var 玩家 Player

	// 選角色
	var 選擇位子 int
	for {
		pos := 遊戲.顯示可選位子()
		連線.WriteJSON(map[string]interface{}{
			"event":    "選擇位子",
			"position": pos,
		})

		_, msg, err := 連線.ReadMessage()
		if err != nil {
			log.Print("選位子錯誤 => ", err)
			遊戲.移除連線(連線)
			return
		}

		err = json.Unmarshal(msg, &選擇位子)
		if err == nil {
			玩家 = 遊戲.玩家們[選擇位子]
			if !玩家.已經被選擇() {
				break
			}
		}
	}

	遊戲.讀寫鎖.RLock()
	是房主 := 連線 == 遊戲.房主
	遊戲.讀寫鎖.RUnlock()
	連線.WriteJSON(map[string]interface{}{
		"event": "你的角色",
		"職業":    玩家.職業(),
		"種族":    玩家.種族(),
		"房主":    是房主,
	})

	玩家.加入(遊戲, 連線)
}

// 開始 開始遊戲
func (遊戲 *Game) 開始() {
	遊戲.讀寫鎖.Lock()
	遊戲.階段 = 開始階段
	遊戲.讀寫鎖.Unlock()

	夜晚 := true
	輪數 := 0

	for {
		輪數++
		if 夜晚 {
			遊戲.天黑請閉眼()
		} else {
			遊戲.天亮請睜眼()
			遊戲.全員請投票()
		}

		遊戲結果 := 遊戲.判斷勝負()
		if 遊戲結果 == 進行中 {
			夜晚 = !夜晚
			continue
		}

		log.Print(遊戲結果)
		for i := range 遊戲.連線池 {
			連線 := 遊戲.連線池[i]
			連線.WriteJSON(遊戲結果)
		}

		遊戲.讀寫鎖.Lock()
		遊戲.階段 = 結束階段
		遊戲.讀寫鎖.Unlock()

		return
	}
}

func (遊戲 *Game) 天黑請閉眼() {
	for i := range 遊戲.連線池 {
		連線 := 遊戲.連線池[i]
		連線.WriteJSON("天黑請閉眼")
	}

	for i := range 遊戲.玩家們 {
		玩家 := 遊戲.玩家們[i]
		玩家.閉眼睛()
	}

	狼人玩家們 := []Skiller{}
	神職玩家們 := []Skiller{}

	for i := range 遊戲.玩家們 {
		玩家 := 遊戲.玩家們[i]
		if 玩家.種族() == 狼職 {
			狼人玩家們 = append(狼人玩家們, 玩家.(Skiller))
		} else if 玩家.職業() != 平民 {
			神職玩家們 = append(神職玩家們, 玩家.(Skiller))
		}
	}

	// 狼人請睜眼
	for i := range 狼人玩家們 {
		狼人 := 狼人玩家們[i]
		狼人.開眼睛()
	}

	// 狼人請殺人
	for i := range 狼人玩家們 {
		狼人 := 狼人玩家們[i]
		狼人.能力()
	}

	// 狼人請閉眼
	for i := range 狼人玩家們 {
		狼人 := 狼人玩家們[i]
		狼人.閉眼睛()
	}

	// 神職請睜眼
	for i := range 神職玩家們 {
		神 := 神職玩家們[i]
		神.開眼睛()
		神.能力()
		神.閉眼睛()
	}
}

func (遊戲 *Game) 天亮請睜眼() {
	for i := range 遊戲.連線池 {
		連線 := 遊戲.連線池[i]
		連線.WriteJSON("天亮請睜眼")
	}

	for i := range 遊戲.玩家們 {
		玩家 := 遊戲.玩家們[i]
		玩家.開眼睛()
	}
}

func (遊戲 *Game) 全員請投票() {

	投票結果 := map[int]int{}
	for i := range 遊戲.玩家們 {
		玩家 := 遊戲.玩家們[i]
		no := 玩家.投票()
		投票結果[玩家.號碼()] = no
	}

	// 顯示給所有玩家看
	for i := range 遊戲.連線池 {
		連線 := 遊戲.連線池[i]
		err := 連線.WriteJSON(投票結果)
		if err != nil {
			遊戲.移除連線(連線)
		}
	}

	// 統計票數
	票數 := map[int]int{}
	for 玩家號碼 := range 投票結果 {
		投給誰 := 投票結果[玩家號碼]
		票數[投給誰]++
	}
}

// 初始設定過 是否初始設定過
func (遊戲 *Game) 初始設定過() bool {
	遊戲.讀寫鎖.Lock()
	設定了 := len(遊戲.玩家們) > 0
	遊戲.讀寫鎖.Unlock()
	return 設定了
}

// 初始設定 初始設定
func (遊戲 *Game) 初始設定(
	選擇角色 map[string]int,
) {
	遊戲.讀寫鎖.Lock()

	隨機可選角色 := []RULE{}
	for 角色, 數量 := range 選擇角色 {
		for i := 0; i < 數量; i++ {
			隨機可選角色 = append(隨機可選角色, RULE(角色))
		}
	}
	隨機可選角色 = 亂數洗牌(隨機可選角色)

	玩家們 := []Player{}
	for i := range 隨機可選角色 {
		switch 隨機可選角色[i] {
		case 平民:
			玩家們 = append(玩家們, NewHuman(i+1))
		case 狼人:
			玩家們 = append(玩家們, NewWolf(i+1))
		}
	}
	遊戲.玩家們 = 玩家們

	遊戲.讀寫鎖.Unlock()
}

// 顯示可選位子 顯示可選位子
func (遊戲 *Game) 顯示可選位子() []int {
	可選位子 := []int{}
	for i := range 遊戲.玩家們 {
		還沒被選擇 := !遊戲.玩家們[i].已經被選擇()
		if 還沒被選擇 {
			可選位子 = append(可選位子, i)
		}
	}
	return 可選位子
}

func (遊戲 *Game) 儲存連線(連線 *websocket.Conn) {
	遊戲.讀寫鎖.Lock()
	if len(遊戲.連線池) == 0 {
		遊戲.房主 = 連線
	}
	遊戲.連線池 = append(遊戲.連線池, 連線)
	遊戲.讀寫鎖.Unlock()
}

func (遊戲 *Game) 移除連線(目前連線 *websocket.Conn) {
	遊戲.讀寫鎖.Lock()
	for i := range 遊戲.連線池 {
		if 遊戲.連線池[i] == 目前連線 {
			遊戲.連線池 = append(遊戲.連線池[:i], 遊戲.連線池[i+1:]...)
		}
	}

	if 目前連線 == 遊戲.房主 {
		if len(遊戲.連線池) > 0 {
			遊戲.房主 = 遊戲.連線池[0]
		} else {
			遊戲.房主 = nil
		}
	}
	遊戲.讀寫鎖.Unlock()
}

func (遊戲 *Game) 是房主(連線 *websocket.Conn) bool {
	是 := 遊戲.房主 == 連線
	return 是
}

func (遊戲 *Game) 判斷勝負() 遊戲結果 {
	神職人數 := 0
	狼職人數 := 0
	平民人數 := 0

	for i := range 遊戲.玩家們 {
		玩家 := 遊戲.玩家們[i]
		玩家還沒出局 := !玩家.出局了()
		if 玩家還沒出局 {
			switch 玩家.種族() {
			case 人質:
				平民人數++
			case 神職:
				神職人數++
			case 狼職:
				狼職人數++
			}
		}
	}

	if 狼職人數 >= 神職人數+平民人數 {
		return 狼勝
	}

	if 平民人數+神職人數+狼職人數 == 0 {
		return 平手
	}

	if 狼職人數 == 0 {
		return 人勝
	}

	return 進行中
}

func (遊戲 *Game) 目前階段() 階段 {
	遊戲.讀寫鎖.RLock()
	目前階段 := 遊戲.階段
	遊戲.讀寫鎖.RUnlock()
	return 目前階段
}
