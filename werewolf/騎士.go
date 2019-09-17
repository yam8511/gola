package werewolf

import (
	"encoding/json"
	"log"
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
}

func (我 *Knight) 種族() GROUP {
	return 神職
}

func (我 *Knight) 職業() RULE {
	return 騎士
}
func (我 *Knight) 發言() {
	log.Println(我.號碼(), " 騎士 開始發言")
	訊息 := <-我.傳話筒
	if string(訊息) == "skill" {
		我.能力()
	}
}

func (我 *Knight) 能力() {
	可指定的玩家號碼 := []int{}
	目前存活玩家們 := 我.遊戲.存活玩家們()
	for i := range 目前存活玩家們 {
		玩家 := 目前存活玩家們[i]
		可指定的玩家號碼 = append(可指定的玩家號碼, 玩家.號碼())
	}

	for {
		我.旁白有話給玩家(map[string]interface{}{
			"event": "請問你要騎誰？",
			"玩家":    可指定的玩家號碼,
		})

		訊息 := <-我.傳話筒
		log.Print(string(訊息))

		玩家號碼 := 0
		err := json.Unmarshal(訊息, &玩家號碼)
		if err != nil {
			continue
		}

		玩家 := 我.遊戲.玩家資料(玩家號碼)
		if 玩家 != nil {
			我.遊戲.殺玩家(騎殺, 玩家)
			if 玩家.種族() == 狼職 {
				我.遊戲.旁白(map[string]interface{}{
					"no":    玩家.號碼(),
					"event": "是狼人！",
				})
				玩家.出局(騎殺)
			} else {
				我.遊戲.旁白(map[string]interface{}{
					"no":    玩家.號碼(),
					"event": "不是狼人！ 以死謝罪",
				})
				我.出局(騎殺)
			}
			return
		}
	}
}

func (我 *Knight) 需要夜晚行動() bool {
	return false
}
