package werewolf

import (
	"encoding/json"
	"strconv"
)

// NewPrettyWolf 建立新PrettyWolf
func NewPrettyWolf(遊戲 *Game, 位子 int) *PrettyWolf {
	wolf := NewWolf(遊戲, 位子)
	return &PrettyWolf{
		Wolf: wolf,
	}
}

// PrettyWolf 狼美人
type PrettyWolf struct {
	*Wolf
	上次睡的玩家 Player
}

func (我 *PrettyWolf) 職業() RULE {
	return 狼美人
}

func (我 *PrettyWolf) 發言(投票發言 bool) bool {

	if 投票發言 {
		我.Human.發言(投票發言)
		return false
	}

	uid := newUID()
	我.遊戲.旁白有話對單個玩家說(我, 傳輸資料{
		UID:     uid,
		Display: "您要發動技能嗎? " + 我.遊戲.提示發言(),
		Action:  等待回應,
		Data: map[string]string{
			"發動✅": "yes",
			"跳過❌": "no",
		},
	}, 0)

	我.等待動作(等待回應, uid)

	return false
}

func (我 *PrettyWolf) 出局(殺法 KILL) {
	我.Human.出局(殺法)
	我.遊戲.判斷勝負(false)
	if 我.遊戲.勝負 == 進行中 {
		我.能力(2)
	}
}

func (我 *PrettyWolf) 能力(i int) (_ Player) {
	switch i {
	case 1:
		if 我.出局了() || !我.已經被選擇() {
			我.遊戲.等一下(random(3) * 2500)
			return
		}

		可睡的玩家號碼 := []int{}
		目前存活玩家們 := 我.遊戲.存活玩家們()
		for i := range 目前存活玩家們 {
			玩家 := 目前存活玩家們[i]
			if 我.號碼() == 玩家.號碼() {
				continue
			}
			可睡的玩家號碼 = append(可睡的玩家號碼, 玩家.號碼())
		}

		uid := newUID()
		我.遊戲.旁白有話對單個玩家說(我, 傳輸資料{
			UID:     uid,
			Display: "請問你要睡的對象是？",
			Action:  選擇玩家,
			Data:    可睡的玩家號碼,
		}, 0)

		for {
			so, err := 我.等待動作(選擇玩家, uid)
			if err != nil {
				return
			}

			玩家號碼 := 0
			err = json.Unmarshal([]byte(so.Reply), &玩家號碼)
			if err != nil {
				continue
			}

			玩家, 存在 := 我.遊戲.玩家資料(玩家號碼)
			if 存在 {
				我.上次睡的玩家 = 玩家
				return
			}
		}

	case 2:
		uid := newUID()
		我.遊戲.旁白(傳輸資料{
			UID:    uid,
			Sound:  strconv.Itoa(我.上次睡的玩家.號碼()) + "號玩家淘汰!" + 我.遊戲.提示點擊(false),
			Action: 等待回應,
		}, 3000)

		for _, 存活玩家 := range 我.遊戲.存活玩家們() {
			存活玩家.等待動作(等待回應, uid)
		}

		我.遊戲.殺玩家(毒殺, 我.上次睡的玩家)
		我.遊戲.玩家出局()
		return

	default:
		return 我.Wolf.能力(0)
	}
}
