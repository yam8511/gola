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

func (我 *WolfKing) 種族() GROUP {
	return 狼職
}

func (我 *WolfKing) 職業() RULE {
	return 狼王
}

func (我 *WolfKing) 需要夜晚行動() bool {
	return false
}

func (我 *WolfKing) 出局(殺法 KILL) {
	我.Human.出局(殺法)
	我.遊戲.判斷勝負(false)
	if 殺法 != 毒殺 && 我.遊戲.勝負 != 進行中 {
		我.能力()
	}
}

func (我 *WolfKing) 能力() (_ Player) {
	if 我.出局了() {
		我.遊戲.旁白(傳輸資料{Sound: "啟動角色技能，請問你要帶走誰？"}, 2000)

		可殺的玩家號碼 := map[string]int{"不帶": -1}
		目前存活玩家們 := 我.遊戲.存活玩家們()
		for i := range 目前存活玩家們 {
			號碼 := 目前存活玩家們[i].號碼()
			可殺的玩家號碼[strconv.Itoa(號碼)] = 號碼
		}

		我.遊戲.旁白有話對單個玩家說(我, 傳輸資料{
			Display: "請問你要帶走誰？",
			Action:  等待回應,
			Data:    可殺的玩家號碼,
		}, 1000)

		for {
			so, err := 我.等待動作(等待回應)
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
				玩家.出局(獵殺)
				我.遊戲.判斷勝負(false)
				return
			}
		}
	}

	return 我.Wolf.能力()
}
