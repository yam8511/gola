package werewolf

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"time"
)

// NewKnight 建立新Knight
func NewKnight(遊戲 *Game, 位子 int) *Knight {
	human := NewHuman(遊戲, 位子)
	return &Knight{
		Human: human,
	}
}

// Knight 玩家
type Knight struct {
	*Human
	發動過技能了 bool
}

func (我 *Knight) 種族() GROUP {
	return 神職
}

func (我 *Knight) 職業() RULE {
	return 騎士
}

func (我 *Knight) 發言(投票發言 bool) bool {

	if 我.發動過技能了 || 投票發言 {
		我.Human.發言(投票發言)
		return false
	}

	uid := newUID()
	我.遊戲.旁白有話對單個玩家說(我, 傳輸資料{
		UID:     uid,
		Display: "您要發動技能嗎?",
		Action:  等待回應,
		Data: map[string]string{
			"發動✅": "yes",
			"跳過❌": "no",
		},
	}, 0)

	so, err := 我.等待動作(等待回應, uid)
	if err == nil && so.Reply == "yes" {
		我.遊戲.旁白(傳輸資料{Sound: strconv.Itoa(我.號碼()) + "號騎士發動技能，請問你要查驗的對象是？"}, 0)

		玩家 := 我.能力()
		if 玩家 != nil && 玩家.種族() == 狼職 {
			return true
		}
	}

	return false
}

func (我 *Knight) 能力() (_ Player) {

	我.發動過技能了 = true

	可指定的玩家號碼 := []int{}
	目前存活玩家們 := 我.遊戲.存活玩家們()
	for _, 玩家 := range 目前存活玩家們 {
		if 我.號碼() == 玩家.號碼() {
			continue
		}
		可指定的玩家號碼 = append(可指定的玩家號碼, 玩家.號碼())
	}

	uid := newUID()
	我.遊戲.旁白有話對單個玩家說(我, 傳輸資料{
		UID:     uid,
		Display: "請問你要騎誰？",
		Action:  選擇玩家,
		Data:    可指定的玩家號碼,
	}, 0)

	玩家號碼 := 0
	for {
		so, err := 我.等待動作(選擇玩家, uid)
		if err != nil {
			rand.Seed(time.Now().UTC().UnixNano())
			rd := rand.Int()
			count := len(可指定的玩家號碼)
			if count == 0 {
				return
			}
			玩家號碼 = 可指定的玩家號碼[rd%count]
		} else {
			err = json.Unmarshal([]byte(so.Reply), &玩家號碼)
			if err != nil {
				continue
			}
		}

		玩家, 存在 := 我.遊戲.玩家資料(玩家號碼)
		if 存在 {
			uid := newUID()
			if 玩家.種族() == 狼職 {
				台詞 := strconv.Itoa(玩家.號碼()) + "號是狼人！大家請點擊確認，即將進入黑夜。"
				我.遊戲.旁白(傳輸資料{
					UID:    uid,
					Sound:  台詞,
					Action: 等待回應,
				}, 3000)
				我.遊戲.殺玩家(騎殺, 玩家)
			} else {
				台詞 := strconv.Itoa(玩家.號碼()) + "號不是狼人！ 騎士以死謝罪。大家請點擊確認，即將進入黑夜。"
				我.遊戲.旁白(傳輸資料{
					UID:    uid,
					Sound:  台詞,
					Action: 等待回應,
				}, 3000)
				我.遊戲.殺玩家(騎殺, 我)
			}

			for _, 存活玩家 := range 目前存活玩家們 {
				存活玩家.等待動作(等待回應, uid)
			}

			return 玩家
		}
	}
}

func (我 *Knight) 需要夜晚行動() bool {
	return false
}
