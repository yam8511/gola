package werewolf

import (
	"encoding/json"
	"log"
	"strconv"
	"time"
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
	// 可殺的玩家號碼 := []int{}
	// 目前存活玩家們 := 我.遊戲.存活玩家們()
	for i := range 目前存活玩家們 {
		玩家 := 目前存活玩家們[i]
		可查看的玩家號碼 = append(可查看的玩家號碼, 玩家.號碼())
	}

	if 我.連線 == nil {
		return
	}

	for {
		err := 我.連線.WriteJSON(map[string]interface{}{
			"event": "請問你要查誰？",
			"玩家":    可查看的玩家號碼,
		})
		if err != nil {
			我.遊戲.移除連線(我.連線)
		}

		訊息 := <-我.傳話筒
		log.Print(string(訊息))

		玩家號碼 := 0
		err = json.Unmarshal(訊息, &玩家號碼)
		if err != nil {
			continue
		}

		玩家 := 我.遊戲.玩家資料(玩家號碼)
		if 玩家.號碼() == 玩家號碼 {
			var s = "是『好人』"

			if 玩家.種族() == 狼職 {
				s = "是『壞人』"
			}
			我.連線.WriteJSON(
				strconv.Itoa(玩家號碼) + "號玩家" + s,
			)
			time.Sleep(3 * time.Second)
			return
		}

	}
}

func (我 *Prophesier) 需要夜晚行動() bool {
	return true
}
