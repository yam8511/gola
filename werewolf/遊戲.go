package werewolf

import (
	"encoding/json"
	"errors"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Game 狼人殺遊戲
type Game struct {
	房主號碼      int
	讀寫鎖       sync.RWMutex
	玩家們       map[string]Player
	淘汰者       []淘汰者
	時間        時間
	階段        階段
	勝負        遊戲結果
	輪數        int
	上一晚狼殺玩家號碼 int
	有神職       bool
	有人質       bool
	獵人場       bool
	角色設定      map[RULE]int
}

func (遊戲 *Game) 加入(連線 *websocket.Conn) {
	進入遊戲 := func(uid string, 玩家 Player) {
		遊戲.旁白有話對單個玩家說(玩家, 傳輸資料{
			Sound:  "你的角色",
			Action: 拿到角色,
			Data: map[string]interface{}{
				"位子":   玩家.號碼(),
				"職業":   玩家.職業(),
				"種族":   玩家.種族(),
				"編號":   uid,
				"角色設定": 遊戲.角色設定,
			}}, 0,
		)

		遊戲.加入玩家(玩家)

		玩家.等待中()
	}

	if 遊戲.目前階段() == 開始階段 {

		for 遊戲.目前階段() == 開始階段 {
			availablePlayer := []int{}
			for _, player := range 遊戲.玩家們 {
				if !player.已經被選擇() {
					availablePlayer = append(availablePlayer, player.號碼())
				}
			}

			err := 遊戲.旁白有話對連線說(連線, 傳輸資料{
				Sound:  "遊戲已經開始，想重新進入遊戲，請輸入玩家號碼",
				Action: 遊戲已開始,
				Data:   availablePlayer,
			}, 0)

			if err != nil {
				return
			}

			so, err := waitSocketBack(連線, 遊戲已開始)
			if err != nil {
				return
			}

			no, err := strconv.Atoi(so.Reply)
			if err != nil {
				continue
			}

			玩家, 存在 := 遊戲.玩家資料(no)
			if 存在 && 玩家.加入(連線) {
				進入遊戲(strconv.Itoa(no), 玩家)
				return
			}
		}
	}

	var 玩家 Player
	var 存在 bool
	var uid string

	// 選角色
	var 選擇位子 int
	for {
		pos := 遊戲.顯示可選位子()
		if len(pos) == 0 {
			return
		}

		err := 連線.WriteJSON(傳輸資料{
			Sound:  "請選擇號碼",
			Action: 選擇號碼,
			Data:   pos,
		})

		if err != nil {
			if 遊戲.房主號碼 == 0 {
				遊戲.重置()
			}
			return
		}

		so, err := waitSocketBack(連線, 選擇號碼)
		if err != nil {
			if 遊戲.房主號碼 == 0 {
				遊戲.重置()
			}
			return
		}

		err = json.Unmarshal([]byte(so.Reply), &選擇位子)
		if err == nil {
			玩家, uid, 存在 = 遊戲.玩家資料WithUID(選擇位子)
			if 存在 && 玩家.加入(連線) {
				break
			}
		}
	}

	進入遊戲(uid, 玩家)
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
		switch {
		case 夜晚:
			遊戲.輪數++
			遊戲.天黑請閉眼()
		case 遊戲.獵人場:
			遊戲.天亮請睜眼()
			遊戲.大家開始發言(nil)
			遊戲.宣布淘汰玩家()
			遊戲.玩家出局()
		default:
			遊戲.天亮請睜眼()
			遊戲.宣布淘汰玩家()
			遊戲.玩家出局()
			遊戲結果 = 遊戲.判斷勝負(false)
			if 遊戲結果 != 進行中 {
				break
			}

			進入天黑 := 遊戲.大家開始發言(nil)
			遊戲.玩家出局()
			遊戲結果 = 遊戲.判斷勝負(false)
			if 遊戲結果 != 進行中 {
				break
			}

			if !進入天黑 {
				遊戲.全員請投票()
			}
		}

		遊戲結果 = 遊戲.判斷勝負(false)
		if 遊戲結果 != 進行中 {
			break
		}
		夜晚 = !夜晚
	}

	sound := "遊戲結束，"
	switch 遊戲結果 {
	case 人勝:
		sound += "好人陣營獲勝"
	case 狼勝:
		sound += "邪惡陣營獲勝"
	case 平手:
		sound += "好人邪惡陣營平手"
	}
	遊戲.旁白(傳輸資料{
		Sound:  sound,
		Action: 遊戲結束,
		Data:   遊戲結果,
	}, 1200)
	runtime.Gosched()

	遊戲.重置()
	return
}

func (遊戲 *Game) 天黑請閉眼() {
	遊戲.時間 = 黑夜
	遊戲.旁白(傳輸資料{
		Sound:  "天黑請閉眼",
		Action: 天黑請閉眼,
	}, 3000)

	for i := range 遊戲.玩家們 {
		玩家 := 遊戲.玩家們[i]
		玩家.閉眼睛()
	}

	狼人玩家們, 神職玩家們 := PickSkiller(遊戲.玩家們)

	// 狼人請睜眼
	遊戲.旁白(傳輸資料{Sound: "狼人請睜眼"}, 3000)
	for i := range 狼人玩家們 {
		狼人 := 狼人玩家們[i]
		狼人.開眼睛()
	}

	// 狼人請殺人
	遊戲.旁白(傳輸資料{Sound: "狼人請殺人"}, 1000)
	if len(狼人玩家們) > 0 {
		狼人們咬殺的對象 := map[Player]bool{}
		wg := sync.WaitGroup{}
		for i := range 狼人玩家們 {
			wg.Add(1)
			go func(i int) {
				狼人 := 狼人玩家們[i]
				if !狼人.出局了() {
					被咬玩家 := 狼人.能力()
					if 被咬玩家 != nil {
						遊戲.讀寫鎖.Lock()
						狼人們咬殺的對象[被咬玩家] = 狼人.已經被選擇()
						遊戲.讀寫鎖.Unlock()
					}
				}
				wg.Done()
			}(i)
			runtime.Gosched()
		}
		wg.Wait()

		// 確認是否有玩家狼
		狼咬數 := map[Player]int{}
		var 狼咬最多的玩家, 玩家狼咬最多的玩家 Player

		for 玩家, 是玩家咬的 := range 狼人們咬殺的對象 {

			狼咬數[玩家]++
			if 狼咬數[玩家] > 狼咬數[狼咬最多的玩家] {
				狼咬最多的玩家 = 玩家
			}

			if 是玩家咬的 {
				if 狼咬數[玩家] > 狼咬數[玩家狼咬最多的玩家] {
					玩家狼咬最多的玩家 = 玩家
				}
			}
		}

		if 玩家狼咬最多的玩家 != nil {
			遊戲.殺玩家(狼殺, 玩家狼咬最多的玩家)
			遊戲.上一晚狼殺玩家號碼 = 玩家狼咬最多的玩家.號碼()
		} else {
			遊戲.殺玩家(狼殺, 狼咬最多的玩家)
			遊戲.上一晚狼殺玩家號碼 = 狼咬最多的玩家.號碼()
		}
	}

	// 狼人請閉眼
	for i := range 狼人玩家們 {
		狼人 := 狼人玩家們[i]
		狼人.閉眼睛()
	}
	遊戲.旁白(傳輸資料{Sound: "狼人請閉眼"}, 3000)

	// 神職請睜眼
	for i := range 神職玩家們 {
		神 := 神職玩家們[i]
		if 神.需要夜晚行動() {
			遊戲.旁白(傳輸資料{Sound: string(神.職業()) + "請睜眼"}, 3000)
			神.開眼睛()
			神.能力()
			神.閉眼睛()
			遊戲.旁白(傳輸資料{Sound: string(神.職業()) + "請閉眼"}, 3000)
		}
	}
}

func (遊戲 *Game) 天亮請睜眼() {
	遊戲.時間 = 白天
	遊戲.旁白(傳輸資料{Sound: "天亮請睜眼", Action: 天亮請睜眼}, 3000)

	for i := range 遊戲.玩家們 {
		玩家 := 遊戲.玩家們[i]
		玩家.開眼睛()
	}
}

func (遊戲 *Game) 宣布淘汰玩家() {
	if len(遊戲.淘汰者) > 0 {
		sort.Slice(遊戲.淘汰者, func(i, j int) bool {
			return 遊戲.淘汰者[i].玩家.號碼() < 遊戲.淘汰者[j].玩家.號碼()
		})

		var 死者名單號碼 []string
		for i := range 遊戲.淘汰者 {
			死者名單號碼 = append(死者名單號碼, strconv.Itoa(遊戲.淘汰者[i].玩家.號碼())+"號")
		}

		uid := newUID()
		遊戲.旁白(傳輸資料{
			UID:    uid,
			Sound:  "昨晚 " + strings.Join(死者名單號碼, ",") + "玩家淘汰!，大家請點擊確認，之後開始發言。",
			Action: 等待回應,
		}, 0)

		for _, 玩家 := range 遊戲.存活玩家們() {
			玩家.等待動作(等待回應, uid)
		}

		if 遊戲.輪數 == 1 {
			for i := range 遊戲.淘汰者 {
				死者 := 遊戲.淘汰者[i]
				遊戲.旁白(傳輸資料{Sound: strconv.Itoa(死者.玩家.號碼()) + "號玩家發表遺言"}, 2000)
				死者.玩家.發表遺言()
			}
		}
	} else {
		遊戲.旁白(傳輸資料{Sound: "昨晚平安夜"}, 3000)
	}
}

func (遊戲 *Game) 大家開始發言(玩家們 []Player) (直接天黑 bool) {
	var 投票發言 bool
	if len(玩家們) == 0 {
		玩家們 = 遊戲.存活玩家們()
	} else {
		投票發言 = true
	}

	for _, 玩家 := range 玩家們 {
		遊戲.旁白(傳輸資料{Sound: strconv.Itoa(玩家.號碼()) + "號玩家開始發言"}, 2000)
		直接天黑 = 玩家.發言(投票發言)
		遊戲結果 := 遊戲.判斷勝負(false)
		if 直接天黑 || 遊戲結果 != 進行中 {
			break
		}
	}

	return
}

func (遊戲 *Game) 全員請投票() {
	投票流程 := func(
		還沒出局的玩家們 []Player,
		可投票玩家號碼 []int,
		投票結果 map[int]int,
		顯示投票文字 string,
	) (
		最高票數 int,
		投票統計結果 map[int]int,
		平票號碼 map[int]int,
	) {
		最高票數 = 0
		投票統計結果 = map[int]int{}
		平票號碼 = map[int]int{}

		投票選項 := map[string]int{
			"棄票": -1,
		}
		for i := range 可投票玩家號碼 {
			號碼 := 可投票玩家號碼[i]
			投票選項[strconv.Itoa(號碼)] = 號碼
		}

		遊戲.旁白(傳輸資料{Sound: 顯示投票文字}, 2000)
		uid := newUID()
		遊戲.旁白(傳輸資料{
			UID:    uid,
			Action: 等待回應,
			Data:   投票選項,
		}, 0)

		顯示用的投票結果 := map[int]interface{}{}
		for i := range 還沒出局的玩家們 {
			玩家 := 還沒出局的玩家們[i]
			投給誰 := 玩家.投票(uid)
			_, 有效投票 := 投票結果[投給誰]
			if 有效投票 {
				投票統計結果[投給誰]++
				顯示用的投票結果[玩家.號碼()] = 投給誰
			} else {
				顯示用的投票結果[玩家.號碼()] = "棄票"
			}
		}

		// 顯示給所有玩家看
		遊戲.旁白(傳輸資料{
			Sound:  "投票結果",
			Action: 顯示投票結果,
			Data:   顯示用的投票結果,
		}, 2000)

		// 統計票數
		for 被投人, 票數 := range 投票統計結果 {
			if 票數 > 最高票數 {
				最高票數 = 票數
				平票號碼 = map[int]int{
					被投人: 最高票數,
				}
			} else if 票數 == 最高票數 {
				平票號碼[被投人] = 最高票數
			}
		}

		return
	}

	可投票玩家號碼 := []int{}
	投票結果 := map[int]int{}
	還沒出局的玩家們 := 遊戲.存活玩家們()
	for i := range 還沒出局的玩家們 {
		玩家號碼 := 還沒出局的玩家們[i].號碼()
		可投票玩家號碼 = append(可投票玩家號碼, 玩家號碼)
		投票結果[玩家號碼] = 0
	}

	最高票數, 投票統計結果, 平票號碼 := 投票流程(
		還沒出局的玩家們,
		可投票玩家號碼,
		投票結果,
		"請投票",
	)

	// 有平票出現，需要投第二輪
	有平票出現 := len(平票號碼) > 1
	if 有平票出現 {

		平票玩家號碼 := []int{}
		for 被投人 := range 平票號碼 {
			平票玩家號碼 = append(平票玩家號碼, 被投人)
		}

		sort.Slice(平票玩家號碼, func(i, j int) bool {
			return 平票玩家號碼[i] < 平票玩家號碼[j]
		})

		平票玩家們 := []Player{}
		平票號碼文字 := []string{}
		for _, 被投人 := range 平票玩家號碼 {
			玩家, 存在 := 遊戲.玩家資料(被投人)
			if 存在 {
				平票玩家們 = append(平票玩家們, 玩家)
				平票號碼文字 = append(平票號碼文字, strconv.Itoa(被投人)+"號")
			}
		}

		遊戲.旁白(傳輸資料{
			Sound: strings.Join(平票號碼文字, ",") + "玩家平票，請重新發言",
		}, 2000)
		遊戲.大家開始發言(平票玩家們)

		最高票數, 投票統計結果, 平票號碼 = 投票流程(
			還沒出局的玩家們,
			平票玩家號碼,
			投票結果,
			"平票，請投票",
		)

		有平票出現 = len(平票號碼) > 1
		if 有平票出現 {
			return
		}
	}

	for 玩家號碼, 票數 := range 投票統計結果 {
		if 票數 == 最高票數 {
			玩家, 存在 := 遊戲.玩家資料(玩家號碼)
			if 存在 {
				遊戲.殺玩家(票殺, 玩家)
				遊戲.旁白(傳輸資料{Sound: strconv.Itoa(玩家號碼) + "號玩家淘汰! 請發表遺言"}, 3000)
				玩家.發表遺言()
			}
			break
		}
	}

	遊戲.玩家出局()
}

func (遊戲 *Game) 玩家資料(號碼 int) (玩家 Player, 存在 bool) {
	遊戲.讀寫鎖.RLock()
	defer 遊戲.讀寫鎖.RUnlock()
	for id := range 遊戲.玩家們 {
		玩家 = 遊戲.玩家們[id]
		if 玩家.號碼() == 號碼 {
			存在 = true
			break
		}
	}

	return
}

func (遊戲 *Game) 玩家資料WithUID(號碼 int) (玩家 Player, uid string, 存在 bool) {
	遊戲.讀寫鎖.RLock()
	defer 遊戲.讀寫鎖.RUnlock()
	for id := range 遊戲.玩家們 {
		玩家 = 遊戲.玩家們[id]
		if 玩家.號碼() == 號碼 {
			uid = id
			存在 = true
			break
		}
	}

	return
}

func (遊戲 *Game) 玩家資料ByUID(uid string) (玩家 Player, 存在 bool) {
	遊戲.讀寫鎖.RLock()
	玩家, 存在 = 遊戲.玩家們[uid]
	遊戲.讀寫鎖.RUnlock()
	return
}

func (遊戲 *Game) 重置() {
	遊戲.讀寫鎖.Lock()
	for i := range 遊戲.玩家們 {
		玩家 := 遊戲.玩家們[i]
		連線 := 玩家.連線()
		if 連線 != nil {
			連線.Close()
		}
	}
	遊戲.玩家們 = map[string]Player{}
	遊戲.淘汰者 = []淘汰者{}
	遊戲.階段 = 準備階段
	遊戲.房主號碼 = 0
	遊戲.勝負 = ""
	遊戲.輪數 = 0
	遊戲.上一晚狼殺玩家號碼 = 0
	遊戲.有神職 = false
	遊戲.有人質 = false
	遊戲.角色設定 = map[RULE]int{}
	遊戲.讀寫鎖.Unlock()
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
	if len(選擇角色) == 0 {
		return
	}

	遊戲.讀寫鎖.Lock()

	遊戲.角色設定 = map[RULE]int{}
	隨機可選角色 := []RULE{}
	for 角色, 數量 := range 選擇角色 {
		if 數量 > 0 {
			遊戲.角色設定[角色] = 數量
			for i := 0; i < 數量; i++ {
				隨機可選角色 = append(隨機可選角色, 角色)
			}
		}
	}
	隨機可選角色 = 亂數洗牌(隨機可選角色)

	遊戲.獵人場 = true
	玩家們 := map[string]Player{}
	for i := range 隨機可選角色 {
		新玩家 := NewPlayer(隨機可選角色[i], 遊戲, i+1)
		if 新玩家 != nil {
			玩家們[strconv.Itoa(i+1)] = 新玩家
			switch 新玩家.種族() {
			case 神職:
				遊戲.有神職 = true
				if 新玩家.職業() != 獵人 {
					遊戲.獵人場 = false
				}
			case 人質:
				遊戲.有人質 = true
				遊戲.獵人場 = false
			}
		}
	}
	遊戲.玩家們 = 玩家們

	遊戲.讀寫鎖.Unlock()
}

func (遊戲 *Game) 旁白(台詞 傳輸資料, 豪秒數 int) {
	for i := range 遊戲.玩家們 {
		玩家 := 遊戲.玩家們[i]
		玩家.傳話給玩家(台詞)
	}

	遊戲.等一下(豪秒數)
}

func (遊戲 *Game) 旁白有話對單個玩家說(玩家 Player, 台詞 傳輸資料, 豪秒數 int) {
	玩家.傳話給玩家(台詞)
	遊戲.等一下(豪秒數)
}

func (遊戲 *Game) 旁白有話對連線說(連線 *websocket.Conn, 台詞 傳輸資料, 豪秒數 int) error {
	if 連線 != nil {
		err := 連線.WriteJSON(台詞)
		if err != nil {
			return errors.New("WebSocket連線斷線了")
		}
	}

	遊戲.等一下(豪秒數)
	return nil
}

func (遊戲 *Game) 等一下(豪秒數 int) {
	if 豪秒數 > 0 {
		time.Sleep(time.Millisecond * time.Duration(豪秒數))
	}
}

func (遊戲 *Game) 存活玩家們() []Player {
	玩家們 := []Player{}
	for i := range 遊戲.玩家們 {
		玩家 := 遊戲.玩家們[i]
		if !玩家.出局了() {
			玩家們 = append(玩家們, 玩家)
		}
	}

	sort.Slice(玩家們, func(i, j int) bool {
		return 玩家們[i].號碼() < 玩家們[j].號碼()
	})

	return 玩家們
}

func (遊戲 *Game) 殺玩家(殺法 KILL, 被殺玩家 Player) {
	if 遊戲.淘汰者 == nil {
		遊戲.淘汰者 = []淘汰者{}
	}

	遊戲.讀寫鎖.Lock()
	遊戲.淘汰者 = append(遊戲.淘汰者, 淘汰者{
		殺法: 殺法,
		玩家: 被殺玩家,
	})
	遊戲.讀寫鎖.Unlock()

	遊戲.判斷勝負(false)
}

func (遊戲 *Game) 救玩家(殺法 KILL) {
	剩餘淘汰者 := make([]淘汰者, len(遊戲.淘汰者))
	copy(剩餘淘汰者, 遊戲.淘汰者)
	for i, 死者 := range 遊戲.淘汰者 {
		if 死者.殺法 == 殺法 {
			剩餘淘汰者 = append(剩餘淘汰者[:i], 剩餘淘汰者[i+1:]...)
		}
	}
	遊戲.淘汰者 = 剩餘淘汰者

	遊戲.判斷勝負(true)
}

func (遊戲 *Game) 尋找淘汰者(殺法 KILL) Player {
	for _, 死者 := range 遊戲.淘汰者 {
		if 死者.殺法 == 殺法 {
			return 死者.玩家
		}
	}

	return nil
}

func (遊戲 *Game) 玩家出局() {
	淘汰者們 := make([]淘汰者, len(遊戲.淘汰者))
	copy(淘汰者們, 遊戲.淘汰者)
	遊戲.淘汰者 = []淘汰者{}
	for _, 死者 := range 淘汰者們 {
		死者.玩家.出局(死者.殺法)
		遊戲.旁白有話對單個玩家說(死者.玩家, 傳輸資料{
			Display: "你已經淘汰，進入觀戰模式",
			Action:  玩家淘汰,
		}, 0)
	}
}

func (遊戲 *Game) 顯示可選位子() []int {
	可選位子 := []int{}

	遊戲.讀寫鎖.Lock()
	for i := range 遊戲.玩家們 {
		還沒被選擇 := !遊戲.玩家們[i].已經被選擇()
		if 還沒被選擇 {
			可選位子 = append(可選位子, 遊戲.玩家們[i].號碼())
		}
	}
	遊戲.讀寫鎖.Unlock()

	sort.Ints(可選位子)
	return 可選位子
}

func (遊戲 *Game) 通知線上人數() {
	線上人數 := []bool{}
	玩家們 := []Player{}

	for i := range 遊戲.玩家們 {
		玩家們 = append(玩家們, 遊戲.玩家們[i])
	}

	sort.Slice(玩家們, func(i, j int) bool {
		return 玩家們[i].號碼() < 玩家們[j].號碼()
	})

	for i := range 玩家們 {
		線上人數 = append(線上人數, 玩家們[i].已經被選擇())
	}

	遊戲.旁白(傳輸資料{
		Action: 更新人數,
		Data:   線上人數,
	}, 0)
}

func (遊戲 *Game) 通知更換房主(玩家 Player) {
	if 遊戲.目前階段() == 準備階段 {
		遊戲.旁白有話對單個玩家說(玩家, 傳輸資料{
			Sound:  "你是房主，隨時可以開始遊戲",
			Action: 更換房主,
			Data: map[string]bool{
				"遊戲開始": 遊戲.目前階段() == 開始階段,
			},
		}, 0)
	} else {
		遊戲.旁白有話對單個玩家說(玩家, 傳輸資料{
			Action: 更換房主,
			Data: map[string]bool{
				"遊戲開始": 遊戲.目前階段() == 開始階段,
			},
		}, 0)
	}
}

func (遊戲 *Game) 加入玩家(玩家 Player) {
	這位玩家是房主 := false
	遊戲.讀寫鎖.Lock()
	if 遊戲.房主號碼 == 0 {
		遊戲.房主號碼 = 玩家.號碼()
		這位玩家是房主 = true
	}
	遊戲.讀寫鎖.Unlock()

	遊戲.通知線上人數()

	if 這位玩家是房主 {
		遊戲.通知更換房主(玩家)
	}
}

func (遊戲 *Game) 踢除玩家(目前玩家 Player) {
	var 房主玩家 Player
	if 遊戲.是房主(目前玩家) {
		遊戲.讀寫鎖.Lock()
		遊戲.房主號碼 = 0
		for i := range 遊戲.玩家們 {
			玩家 := 遊戲.玩家們[i]
			if 玩家.已經被選擇() {
				遊戲.房主號碼 = 玩家.號碼()
				房主玩家 = 玩家
				break
			}
		}
		遊戲.讀寫鎖.Unlock()

		if 遊戲.房主號碼 == 0 {
			遊戲.重置()
		} else {
			遊戲.通知更換房主(房主玩家)
		}
	}
	遊戲.通知線上人數()
}

func (遊戲 *Game) 是房主(玩家 Player) bool {
	遊戲.讀寫鎖.Lock()
	是 := 遊戲.房主號碼 == 玩家.號碼()
	遊戲.讀寫鎖.Unlock()
	return 是
}

func (遊戲 *Game) 判斷勝負(有救人 bool) 遊戲結果 {
	if 遊戲.勝負 != 進行中 && !有救人 {
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

	for i := range 遊戲.淘汰者 {
		玩家 := 遊戲.淘汰者[i].玩家
		switch 玩家.種族() {
		case 人質:
			平民人數--
		case 神職:
			神職人數--
		case 狼職:
			狼職人數--
		}
	}

	if 狼職人數 > 神職人數+平民人數 {
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

	if 平民人數 == 0 && 遊戲.有人質 {
		遊戲.勝負 = 狼勝
		return 狼勝
	}

	if 神職人數 == 0 && 遊戲.有神職 {
		遊戲.勝負 = 狼勝
		return 狼勝
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
