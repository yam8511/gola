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
	抓到狼 bool
}

func (我 *Knight) 種族() GROUP {
	return 神職
}

func (我 *Knight) 職業() RULE {
	return 騎士
}
func (我 *Knight) 發言() bool {
	我.遊戲.旁白有話對單個玩家說(我, 傳輸資料{
		Display: "您有技能可以發動",
		Action:  通知玩家有技能,
		Data:    "skill",
	})

	so, err := waitChannelBack(我.傳話筒, 通知玩家有技能)
	if err != nil {
		return false
	}

	if so.Reply == "skill" {
		我.遊戲.旁白(傳輸資料{Sound: "騎士發動技能"})

		我.能力()
		if 我.抓到狼 {
			return true
		}
	}

	return false
}

func (我 *Knight) 能力() (_ Player) {
	可指定的玩家號碼 := []int{}
	目前存活玩家們 := 我.遊戲.存活玩家們()
	for i := range 目前存活玩家們 {
		玩家 := 目前存活玩家們[i]
		可指定的玩家號碼 = append(可指定的玩家號碼, 玩家.號碼())
	}

	我.遊戲.旁白(傳輸資料{
		Sound:  "請問你要騎誰？",
		Action: 選擇玩家,
		Data:   可指定的玩家號碼,
	})

	玩家號碼 := 0
	for {
		so, err := waitChannelBack(我.傳話筒, 選擇玩家)
		if err != nil {
			rand.Seed(time.Now().UTC().UnixNano())
			rd := rand.Int()
			玩家號碼 = 可指定的玩家號碼[rd%len(可指定的玩家號碼)]
		} else {
			err = json.Unmarshal([]byte(so.Reply), &玩家號碼)
			if err != nil {
				continue
			}
		}

		玩家, 存在 := 我.遊戲.玩家資料(玩家號碼)
		if 存在 {
			我.遊戲.殺玩家(騎殺, 玩家)
			if 玩家.種族() == 狼職 {
				台詞 := strconv.Itoa(玩家.號碼()) + "是狼人！"
				我.遊戲.旁白(傳輸資料{Sound: 台詞})
				我.遊戲.殺玩家(騎殺, 玩家)
				我.抓到狼 = true
			} else {
				台詞 := strconv.Itoa(玩家.號碼()) + "不是狼人！ 騎士以死謝罪"
				我.遊戲.旁白(傳輸資料{Sound: 台詞})
				我.遊戲.殺玩家(騎殺, 我)
			}
			return 玩家
		}
	}
}

func (我 *Knight) 需要夜晚行動() bool {
	return false
}
