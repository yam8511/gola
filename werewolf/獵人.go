package werewolf

import (
	"encoding/json"
	"log"
)

// NewHunter 建立新Hunter
func NewHunter(遊戲 *Game, 位子 int) *Hunter {
	human := NewHuman(遊戲, 位子)
	return &Hunter{
		Human: human,
	}
}

// Hunter 玩家
type Hunter struct {
	*Human
}

func (我 *Hunter) 種族() GROUP {
	return 神職
}

func (我 *Hunter) 職業() RULE {
	return 獵人
}

func (我 *Hunter) 能力() {

	我.遊戲.旁白("獵人玩家已死亡，請在臨走前開最後一槍...請選擇！！！")

	可殺的玩家號碼 := []int{}
	目前存活玩家們 := 我.遊戲.存活玩家們()
	for i := range 目前存活玩家們 {
		玩家 := 目前存活玩家們[i]
		可殺的玩家號碼 = append(可殺的玩家號碼, 玩家.號碼())
	}

	if 我.連線 == nil {
		return
	}

	for {
		err := 我.連線.WriteJSON(map[string]interface{}{
			"event": "請問你要獵殺誰？",
			"玩家":    可殺的玩家號碼,
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

		if 玩家號碼 != 0 {
			for i := range 我.遊戲.玩家們 {
				玩家 := 我.遊戲.玩家們[i]
				if 玩家.號碼() == 玩家號碼 {
					我.遊戲.殺玩家(獵殺, 玩家)
					return
				}
			}
		}
	}
}
