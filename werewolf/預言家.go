package werewolf

import (
	"encoding/json"
	"strconv"
)

// NewProphesier 建立新Prophesier
func NewProphesier(遊戲 *Game, 位子 int) *Prophesier {
	human := NewHuman(遊戲, 位子)
	return &Prophesier{
		Human: human,
	}
}

// Prophesier 玩家
type Prophesier struct {
	*Human
}

func (我 *Prophesier) 種族() GROUP {
	return 神職
}

func (我 *Prophesier) 職業() RULE {
	return 預言家
}

func (我 *Prophesier) 能力() {

	可查看的玩家號碼 := []int{}

	目前存活玩家們 := 我.遊戲.存活玩家們()
	for i := range 目前存活玩家們 {
		玩家 := 目前存活玩家們[i]
		可查看的玩家號碼 = append(可查看的玩家號碼, 玩家.號碼())
	}

	if 我.連線() == nil {
		return
	}

	for {
		我.遊戲.旁白(傳輸資料{
			Sound:  "請問你要查誰？",
			Action: 選擇玩家,
			Data:   可查看的玩家號碼,
		})
		訊息 := <-我.傳話筒

		玩家號碼 := 0
		err := json.Unmarshal(訊息, &玩家號碼)
		if err != nil {
			continue
		}

		玩家 := 我.遊戲.玩家資料(玩家號碼)
		if 玩家 != nil {
			var s = "是『好人』"

			if 玩家.種族() == 狼職 && 玩家.職業() != 雪狼 {
				s = "是『狼人』"
			}

			我.遊戲.旁白有話對單個玩家說(我, 傳輸資料{
				Display: strconv.Itoa(玩家號碼) + s,
			})
			<-我.傳話筒
			return
		}
	}
}

func (我 *Prophesier) 需要夜晚行動() bool {
	return true
}
