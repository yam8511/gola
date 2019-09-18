package werewolf

import (
	"encoding/json"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Game sera找明俊玩遊戲
type Game struct {
	房主號碼  int
	讀寫鎖   sync.RWMutex
	玩家們   []Player
	階段    階段
	通訊    chan int
	夜晚淘汰者 map[KILL]Player
	勝負    遊戲結果
	輪數    int
}

func (遊戲 *Game) 加入(連線 *websocket.Conn) {
	遊戲.讀寫鎖.RLock()
	遊戲已經開始 := 遊戲.階段 == 開始階段
	遊戲.讀寫鎖.RUnlock()
	if 遊戲已經開始 {
		連線.WriteJSON(傳輸資料{
			Sound:  "遊戲已經開始",
			Action: 遊戲已開始,
		})
		return
	}

	var 玩家 Player

	// 選角色
	var 選擇位子 int
	for {
		pos := 遊戲.顯示可選位子()
		連線.WriteJSON(傳輸資料{
			Sound:  "請選擇號碼",
			Action: 選擇玩家,
			Data:   pos,
		})

		_, msg, err := 連線.ReadMessage()
		if err != nil {
			return
		}

		err = json.Unmarshal(msg, &選擇位子)
		if err == nil {
			玩家 = 遊戲.玩家們[選擇位子-1]
			if !玩家.已經被選擇() {
				break
			}
		}
	}

	遊戲.加入玩家(玩家)

	遊戲.旁白有話對連線說(連線, 傳輸資料{
		Sound:  "你的角色",
		Action: 拿到角色,
		Data: map[string]interface{}{
			"位子": 玩家.號碼(),
			"職業": 玩家.職業(),
			"種族": 玩家.種族(),
			"房主": 遊戲.是房主(玩家),
		}},
	)
	玩家.加入(連線)
}

func (遊戲 *Game) 開始() {
	遊戲.讀寫鎖.Lock()
	遊戲.階段 = 開始階段
	遊戲.勝負 = 進行中
	遊戲.讀寫鎖.Unlock()

	夜晚 := true
	遊戲.輪數 = 0

	var 遊戲結果 遊戲結果
	for {
		if 夜晚 {
			遊戲.輪數++
			遊戲.天黑請閉眼()
		} else {
			遊戲.天亮請睜眼()
			遊戲結果 = 遊戲.判斷勝負()
			if 遊戲結果 != 進行中 {
				break
			}

			遊戲.大家開始發言()
			遊戲結果 = 遊戲.判斷勝負()
			if 遊戲結果 != 進行中 {
				break
			}

			遊戲.全員請投票()
		}

		遊戲結果 = 遊戲.判斷勝負()
		if 遊戲結果 != 進行中 {
			break
		}
		夜晚 = !夜晚
	}

	遊戲.旁白(傳輸資料{
		Sound:  "遊戲結束",
		Action: 遊戲結束,
		Data:   遊戲結果,
	})
	遊戲.讀寫鎖.Lock()
	遊戲.階段 = 準備階段
	遊戲.讀寫鎖.Unlock()

	return

}

func (遊戲 *Game) 天黑請閉眼() {
	遊戲.旁白(傳輸資料{
		Sound:  "天黑請閉眼",
		Action: 天黑請閉眼,
	})

	for i := range 遊戲.玩家們 {
		玩家 := 遊戲.玩家們[i]
		玩家.閉眼睛()
	}

	狼人玩家們, 神職玩家們 := PickSkiller(遊戲.玩家們)

	// 狼人請睜眼
	遊戲.旁白(傳輸資料{Sound: "狼人請睜眼"})
	for i := range 狼人玩家們 {
		狼人 := 狼人玩家們[i]
		狼人.開眼睛()
	}

	// 狼人請殺人
	遊戲.旁白(傳輸資料{Sound: "狼人請殺人"})
	for i := range 狼人玩家們 {
		狼人 := 狼人玩家們[i]
		狼人.能力()
	}

	// 狼人請閉眼
	for i := range 狼人玩家們 {
		狼人 := 狼人玩家們[i]
		狼人.閉眼睛()
	}
	遊戲.旁白(傳輸資料{Sound: "狼人請閉眼"})

	// 神職請睜眼
	for i := range 神職玩家們 {
		神 := 神職玩家們[i]
		if 神.需要夜晚行動() {
			遊戲.旁白(傳輸資料{Sound: string(神.職業()) + "請睜眼"})
			神.開眼睛()
			神.能力()
			神.閉眼睛()
			遊戲.旁白(傳輸資料{Sound: string(神.職業()) + "請閉眼"})
		}
	}
}

func (遊戲 *Game) 天亮請睜眼() {
	遊戲.旁白(傳輸資料{Sound: "天亮請睜眼", Action: 天亮請睜眼})

	for i := range 遊戲.玩家們 {
		玩家 := 遊戲.玩家們[i]
		玩家.開眼睛()
	}

	if len(遊戲.夜晚淘汰者) > 0 {
		// 公布淘汰者
		死者名單 := []string{}

		for 殺法 := range 遊戲.夜晚淘汰者 {
			死者名單 = append(死者名單, strconv.Itoa(遊戲.夜晚淘汰者[殺法].號碼()))
		}
		遊戲.旁白(傳輸資料{Sound: "昨晚 " + strings.Join(死者名單, ",") + " 淘汰!"})

		if 遊戲.輪數 == 1 {
			for 殺法 := range 遊戲.夜晚淘汰者 {
				死者 := 遊戲.夜晚淘汰者[殺法]
				遊戲.旁白(傳輸資料{
					Sound:  strconv.Itoa(死者.號碼()) + "號玩家發表遺言",
					Action: 等待回應,
				})
				死者.發表遺言()
			}
		}

		遊戲.玩家出局()
	} else {
		遊戲.旁白(傳輸資料{Sound: "昨晚平安夜"})
	}
}

func (遊戲 *Game) 大家開始發言() {

	存活玩家們 := 遊戲.存活玩家們()
	for i := range 存活玩家們 {
		玩家 := 存活玩家們[i]
		遊戲.旁白(傳輸資料{Sound: strconv.Itoa(玩家.號碼()) + "號玩家開始發言"})
		中斷發話 := 玩家.發言()
		遊戲結果 := 遊戲.判斷勝負()
		if 中斷發話 || 遊戲結果 != 進行中 {
			break
		}
	}

	遊戲.玩家出局()
}

func (遊戲 *Game) 全員請投票() {
	可投票玩家號碼 := map[int]int{}
	還沒出局的玩家們 := 遊戲.存活玩家們()
	for i := range 還沒出局的玩家們 {
		玩家 := 還沒出局的玩家們[i]
		可投票玩家號碼[玩家.號碼()] = 0
	}

	遊戲.旁白(傳輸資料{
		Sound:  "請投票",
		Action: 選擇玩家,
		Data:   可投票玩家號碼,
	})

	投票結果 := map[int]int{}
	for i := range 還沒出局的玩家們 {
		玩家 := 還沒出局的玩家們[i]
		投給誰 := 玩家.投票()
		_, 有效投票 := 可投票玩家號碼[投給誰]
		if 有效投票 {
			投票結果[玩家.號碼()] = 投給誰
		}
	}

	// 顯示給所有玩家看
	遊戲.旁白(傳輸資料{
		Sound:  "投票結果",
		Action: 顯示投票結果,
		Data:   投票結果,
	})

	// 統計票數
	最高票數 := 0
	平票號碼 := map[int]int{}
	for 玩家號碼 := range 投票結果 {
		投給誰 := 投票結果[玩家號碼]
		可投票玩家號碼[投給誰]++
		if 可投票玩家號碼[投給誰] > 最高票數 {
			最高票數 = 可投票玩家號碼[投給誰]
			平票號碼 = map[int]int{
				投給誰: 最高票數,
			}
		} else if 可投票玩家號碼[投給誰] == 最高票數 {
			平票號碼[投給誰] = 最高票數
		}
	}

	// 有平票出現，需要投第二輪
	有平票出現 := len(平票號碼) > 1
	if 有平票出現 {
		遊戲.旁白(傳輸資料{
			Sound:  "請投票",
			Action: 選擇玩家,
			Data:   平票號碼,
		})

		投票結果 := map[int]int{}
		for i := range 還沒出局的玩家們 {
			玩家 := 還沒出局的玩家們[i]
			投給誰 := 玩家.投票()
			_, 有效投票 := 平票號碼[投給誰]
			if 有效投票 {
				投票結果[玩家.號碼()] = 投給誰
			}
		}

		// 顯示給所有玩家看
		遊戲.旁白(傳輸資料{
			Sound:  "投票結果",
			Action: 顯示投票結果,
			Data:   投票結果,
		})

		// 統計票數
		最高票數 = 0
		可投票玩家號碼 = map[int]int{}
		平票號碼 = map[int]int{}
		for 玩家號碼 := range 投票結果 {
			投給誰 := 投票結果[玩家號碼]
			可投票玩家號碼[投給誰]++
			if 可投票玩家號碼[投給誰] > 最高票數 {
				最高票數 = 可投票玩家號碼[投給誰]
				平票號碼 = map[int]int{
					投給誰: 最高票數,
				}
			} else if 可投票玩家號碼[投給誰] == 最高票數 {
				平票號碼[投給誰] = 最高票數
			}
		}

		有平票出現 = len(平票號碼) > 1
		if 有平票出現 {
			return
		}
	}

	for 玩家號碼, 票數 := range 可投票玩家號碼 {
		if 票數 == 最高票數 {
			玩家 := 遊戲.玩家資料(玩家號碼)
			if 玩家 != nil {
				遊戲.殺玩家(票殺, 玩家)
				遊戲.旁白(傳輸資料{
					Sound:  strconv.Itoa(玩家號碼) + " 淘汰! 請發表遺言",
					Action: 等待回應,
				})
				玩家.發表遺言()
			}
			return
		}
	}

	遊戲.玩家出局()
}

func (遊戲 *Game) 玩家資料(號碼 int) Player {
	for i := range 遊戲.玩家們 {
		玩家 := 遊戲.玩家們[i]
		if 玩家.號碼() == 號碼 {
			return 玩家
		}
	}
	return nil
}

func (遊戲 *Game) 初始設定過() bool {
	遊戲.讀寫鎖.Lock()
	設定了 := len(遊戲.玩家們) > 0
	遊戲.讀寫鎖.Unlock()
	return 設定了
}

func (遊戲 *Game) 初始設定(
	選擇角色 map[RULE]int,
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
		新玩家 := NewPlayer(隨機可選角色[i], 遊戲, i+1)
		if 新玩家 != nil {
			玩家們 = append(玩家們, 新玩家)
		}
	}
	遊戲.玩家們 = 玩家們

	遊戲.讀寫鎖.Unlock()
}

func (遊戲 *Game) 旁白(台詞 傳輸資料) {
	for i := range 遊戲.玩家們 {
		玩家 := 遊戲.玩家們[i]
		連線 := 玩家.連線()
		if 連線 != nil {
			err := 連線.WriteJSON(台詞)
			if err != nil {
				玩家.退出()
			}
		}
	}

	遊戲.等一下()
}

func (遊戲 *Game) 旁白有話對單個玩家說(玩家 Player, 台詞 傳輸資料) {

	連線 := 玩家.連線()
	if 連線 != nil {
		err := 遊戲.旁白有話對連線說(連線, 台詞)
		if err != nil {
			玩家.退出()
		}
	}

	遊戲.等一下()
}

func (遊戲 *Game) 旁白有話對連線說(連線 *websocket.Conn, 台詞 傳輸資料) error {
	if 連線 != nil {
		err := 連線.WriteJSON(台詞)
		if err != nil {
			return err
		}
	}

	遊戲.等一下()
	return nil
}

func (遊戲 *Game) 等一下() {
	time.Sleep(time.Millisecond * 500)
}

func (遊戲 *Game) 存活玩家們() []Player {
	還沒出局的玩家 := []Player{}
	for i := range 遊戲.玩家們 {
		玩家 := 遊戲.玩家們[i]
		if !玩家.出局了() {
			還沒出局的玩家 = append(還沒出局的玩家, 玩家)
		}
	}

	return 還沒出局的玩家
}

func (遊戲 *Game) 殺玩家(殺法 KILL, 被殺玩家 Player) {
	if 遊戲.夜晚淘汰者 == nil {
		遊戲.夜晚淘汰者 = map[KILL]Player{}
	}
	遊戲.夜晚淘汰者[殺法] = 被殺玩家

	遊戲.判斷勝負()
}

func (遊戲 *Game) 玩家出局() {
	for 殺法 := range 遊戲.夜晚淘汰者 {
		遊戲.夜晚淘汰者[殺法].出局(殺法)
	}
	遊戲.夜晚淘汰者 = map[KILL]Player{}
}

func (遊戲 *Game) 顯示可選位子() []int {
	可選位子 := []int{}

	for i := range 遊戲.玩家們 {
		還沒被選擇 := !遊戲.玩家們[i].已經被選擇()
		if 還沒被選擇 {
			可選位子 = append(可選位子, 遊戲.玩家們[i].號碼())
		}
	}
	return 可選位子
}

func (遊戲 *Game) 加入玩家(玩家 Player) {
	遊戲.讀寫鎖.Lock()
	if 遊戲.房主號碼 == 0 {
		遊戲.房主號碼 = 玩家.號碼()
	}
	遊戲.讀寫鎖.Unlock()
}

func (遊戲 *Game) 踢除玩家(目前玩家 Player) {
	遊戲.讀寫鎖.Lock()
	if 目前玩家.號碼() == 遊戲.房主號碼 {
		if len(遊戲.玩家們) > 0 {
			for i := range 遊戲.玩家們 {
				玩家 := 遊戲.玩家們[i]
				if 玩家.已經被選擇() {
					遊戲.房主號碼 = 玩家.號碼()
					break
				}
				遊戲.房主號碼 = 0
			}
		} else {
			遊戲.房主號碼 = 0
		}
	}
	遊戲.讀寫鎖.Unlock()
}

func (遊戲 *Game) 是房主(玩家 Player) bool {
	遊戲.讀寫鎖.Lock()
	是 := 遊戲.房主號碼 == 玩家.號碼()
	遊戲.讀寫鎖.Unlock()
	return 是
}

func (遊戲 *Game) 判斷勝負() 遊戲結果 {
	if 遊戲.勝負 != 進行中 {
		return 遊戲.勝負
	}

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

	for 殺法 := range 遊戲.夜晚淘汰者 {
		玩家 := 遊戲.夜晚淘汰者[殺法]
		switch 玩家.種族() {
		case 人質:
			平民人數--
		case 神職:
			神職人數--
		case 狼職:
			狼職人數--
		}
	}

	if 狼職人數 >= 神職人數+平民人數 {
		遊戲.勝負 = 狼勝
		return 狼勝
	}

	if 平民人數+神職人數+狼職人數 == 0 {
		遊戲.勝負 = 平手
		return 平手
	}

	if 狼職人數 == 0 {
		遊戲.勝負 = 人勝
		return 人勝
	}

	遊戲.勝負 = 進行中
	return 進行中
}

func (遊戲 *Game) 目前階段() 階段 {
	遊戲.讀寫鎖.RLock()
	目前階段 := 遊戲.階段
	遊戲.讀寫鎖.RUnlock()
	return 目前階段
}
