package werewolf

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"time"
)

// NewWolf 建立新Wolf
func NewWolf(遊戲 *Game, 位子 int) *Wolf {
	human := NewHuman(遊戲, 位子)
	return &Wolf{
		Human: human,
	}
}

// Wolf 玩家
type Wolf struct {
	*Human
}

func (我 *Wolf) 種族() GROUP {
	return 狼職
}

func (我 *Wolf) 職業() RULE {
	return 狼人
}

func (我 *Wolf) 發言(投票發言 bool) bool {

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

func (我 *Wolf) 能力() (_ Player) {
	可殺的玩家號碼 := []int{}
	目前存活玩家們 := 我.遊戲.存活玩家們()
	for i := range 目前存活玩家們 {
		玩家 := 目前存活玩家們[i]
		if 我.遊戲.上一晚狼殺玩家號碼 == 玩家.號碼() {
			continue
		}
		可殺的玩家號碼 = append(可殺的玩家號碼, 玩家.號碼())
	}

	uid := newUID()
	我.遊戲.旁白有話對單個玩家說(我, 傳輸資料{
		UID:     uid,
		Display: "請問你要殺誰？",
		Action:  選擇玩家,
		Data:    可殺的玩家號碼,
	}, 0)

	for {
		so, err := 我.等待動作(選擇玩家, uid)
		if err != nil {
			rand.Seed(time.Now().UTC().UnixNano())
			rd := rand.Int()
			count := len(可殺的玩家號碼)
			if count == 0 {
				return
			}

			玩家號碼 := 可殺的玩家號碼[rd%count]
			玩家, 存在 := 我.遊戲.玩家資料(玩家號碼)
			if 存在 {
				return 玩家
			}

			continue
		}

		玩家號碼 := 0
		err = json.Unmarshal([]byte(so.Reply), &玩家號碼)
		if err != nil {
			continue
		}

		玩家, 存在 := 我.遊戲.玩家資料(玩家號碼)
		if 存在 {
			// 我.遊戲.殺玩家(狼殺, 玩家)
			return 玩家
		}
	}
}

func (我 *Wolf) 需要夜晚行動() bool {
	return true
}
