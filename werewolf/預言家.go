package werewolf

import (
	"encoding/json"
	"strconv"
)

// NewSeer 建立新Seer
func NewSeer(遊戲 *Game, 位子 int) *Seer {
	human := NewHuman(遊戲, 位子)
	return &Seer{
		Human: human,
	}
}

// Seer 玩家
type Seer struct {
	*Human
}

func (我 *Seer) 種族() GROUP {
	return 神職
}

func (我 *Seer) 職業() RULE {
	return 預言家
}

func (我 *Seer) 需要夜晚行動() bool {
	return true
}

func (我 *Seer) 能力(i int) (_ Player) {

	我.遊戲.旁白(傳輸資料{Sound: "請問你要查驗的對象是？"}, 1500)

	if 我.出局了() || !我.已經被選擇() {
		我.遊戲.等一下(random(3) * 2500)
		return
	}

	可查看的玩家號碼 := []int{}

	目前存活玩家們 := 我.遊戲.存活玩家們()
	for i := range 目前存活玩家們 {
		玩家 := 目前存活玩家們[i]
		if 我.號碼() == 玩家.號碼() {
			continue
		}
		可查看的玩家號碼 = append(可查看的玩家號碼, 玩家.號碼())
	}

	uid := newUID()
	我.遊戲.旁白有話對單個玩家說(我, 傳輸資料{
		UID:     uid,
		Display: "請問你要查驗的對象是？",
		Action:  選擇玩家,
		Data:    可查看的玩家號碼,
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
			var s = "是『😎好人』👍"

			if 玩家.種族() == 狼職 && 玩家.職業() != 雪狼 {
				s = "是『🐺狼人』😱"
			}

			uid := newUID()
			我.遊戲.旁白有話對單個玩家說(我, 傳輸資料{
				UID:     uid,
				Display: strconv.Itoa(玩家號碼) + s,
				Action:  等待回應,
			}, 0)

			我.等待動作(等待回應, uid)
			return
		}
	}
}
