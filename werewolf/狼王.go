package werewolf

import (
	"encoding/json"
	"strconv"
)

// NewWolfKing 建立新WolfKing
func NewWolfKing(遊戲 *Game, 位子 int) *WolfKing {
	wolf := NewWolf(遊戲, 位子)
	return &WolfKing{
		Wolf: wolf,
	}
}

// WolfKing 狼王
type WolfKing struct {
	*Wolf
}

func (我 *WolfKing) 職業() RULE {
	return 狼王
}

func (我 *WolfKing) 發言(投票發言 bool) bool {

	if 投票發言 {
		我.Human.發言(投票發言)
		return false
	}

	uid := newUID()
	我.遊戲.旁白有話對單個玩家說(我, 傳輸資料{
		UID:     uid,
		Display: "您要發動技能嗎? (狼人發動可自爆，騎士發動可查驗)",
		Action:  等待回應,
		Data: map[string]string{
			"發動✅": "yes",
			"跳過❌": "no",
		},
	}, 0)

	so, err := 我.等待動作(等待回應, uid)
	if err == nil && so.Reply == "yes" {

		uid := newUID()
		我.遊戲.旁白(傳輸資料{
			UID:    uid,
			Sound:  strconv.Itoa(我.號碼()) + "號玩家自爆。請點擊確認，即將進入黑夜。",
			Action: 等待回應,
		}, 2000)

		for _, 存活玩家 := range 我.遊戲.存活玩家們() {
			存活玩家.等待動作(等待回應, uid)
		}

		我.遊戲.殺玩家(自爆, 我)

		return true
	}

	return false
}

func (我 *WolfKing) 出局(殺法 KILL) {
	我.Human.出局(殺法)
	我.遊戲.判斷勝負(false)
	if 殺法 != 毒殺 && 我.遊戲.勝負 == 進行中 {
		我.能力()
	}
}

func (我 *WolfKing) 能力() (_ Player) {
	if 我.出局了() {
		我.遊戲.旁白(傳輸資料{Sound: strconv.Itoa(我.號碼()) + "號玩家啟動角色技能，請問你要帶走誰？"}, 3000)

		可殺的玩家號碼 := map[string]int{"不帶": -1}
		目前存活玩家們 := 我.遊戲.存活玩家們()
		for i := range 目前存活玩家們 {
			號碼 := 目前存活玩家們[i].號碼()
			可殺的玩家號碼[strconv.Itoa(號碼)] = 號碼
		}

		uid := newUID()
		我.遊戲.旁白有話對單個玩家說(我, 傳輸資料{
			UID:     uid,
			Display: "請問你要帶走誰？",
			Action:  等待回應,
			Data:    可殺的玩家號碼,
		}, 1000)

		for {
			so, err := 我.等待動作(等待回應, uid)
			if err != nil {
				return
			}

			玩家號碼 := 0
			err = json.Unmarshal([]byte(so.Reply), &玩家號碼)
			if err != nil {
				continue
			}

			if 玩家號碼 == -1 {
				return
			}

			玩家, 存在 := 我.遊戲.玩家資料(玩家號碼)
			if 存在 {
				uid := newUID()
				我.遊戲.旁白(傳輸資料{
					UID: uid,
					Sound: strconv.Itoa(玩家.號碼()) + "號玩家淘汰!請大家點擊確認，即將進入黑夜。",
					Action: 等待回應,
				}, 3000)
				我.遊戲.殺玩家(狼殺, 玩家)
				我.遊戲.判斷勝負(false)

				for _, 存活玩家 := range 目前存活玩家們 {
					存活玩家.等待動作(等待回應, uid)
				}

				return
			}
		}
	}

	return 我.Wolf.能力()
}
