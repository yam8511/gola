package werewolf

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// Witch 玩家
type Witch struct {
	*Human
	有毒藥 bool
	有解藥 bool
}

// NewWitch 建立新女巫
func NewWitch(遊戲 *Game, 位子 int) *Witch {
	human := NewHuman(遊戲, 位子)
	return &Witch{
		Human: human,
		有毒藥:   true,
		有解藥:   true,
	}
}

func (我 *Witch) 種族() GROUP {
	return 神職
}

func (我 *Witch) 職業() RULE {
	return 女巫
}

func (我 *Witch) 能力() (_ Player) {

	有使用過藥了嗎 := false

	if 我.有解藥 {
		我.遊戲.旁白(傳輸資料{Sound: "他被殺了，你要救他嗎？"}, 1500)

		被狼殺玩家 := 我.遊戲.尋找淘汰者(狼殺)

		// 女巫是否救人
		if 被狼殺玩家 != nil && !我.出局了() {

			可以救 := !(我.遊戲.輪數 >= 2 && 被狼殺玩家.號碼() == 我.號碼())
			救的選項 := map[string]string{
				"不救": "no",
			}
			顯示文字 := fmt.Sprintf("你要救『%d號玩家』嗎？", 被狼殺玩家.號碼())
			if 可以救 {
				救的選項["救"] = "yes"
			} else {
				顯示文字 += "  (不能救自己唷哈哈哈笨蛋)"
			}

			uid := newUID()
			我.遊戲.旁白有話對單個玩家說(我, 傳輸資料{
				UID:     uid,
				Display: 顯示文字,
				Action:  等待回應,
				Data:    救的選項,
			}, 0)

			是否使用解藥, err := 我.等待動作(等待回應, uid)

			if err != nil || 是否使用解藥.Reply == "yes" {
				我.遊戲.救玩家(狼殺)
				我.有解藥 = false
				有使用過藥了嗎 = true
			}
		} else {
			我.遊戲.等一下(random(3) * 2500)
		}
	}

	我.遊戲.旁白(傳輸資料{Sound: "你要使用毒藥嗎？你要毒誰呢？"}, 1500)

	if !有使用過藥了嗎 && 我.有毒藥 && !我.出局了() {

		uid := uuid.New().String()
		我.遊戲.旁白有話對單個玩家說(我, 傳輸資料{
			UID:     uid,
			Display: "你要使用毒藥嗎？",
			Action:  等待回應,
			Data: map[string]interface{}{
				"要":  "yes",
				"不要": "no",
			},
		}, 0)

		//如果是第二晚之後被殺的是女巫，則不能使用解藥
		是否使用毒藥, err := 我.等待動作(等待回應, uid)
		if err != nil {
			return
		}

		if 是否使用毒藥.Reply == "yes" {

			我.有毒藥 = false
			可毒的玩家號碼 := []int{}
			目前存活玩家們 := 我.遊戲.存活玩家們()

			for i := range 目前存活玩家們 {
				玩家 := 目前存活玩家們[i]
				if 我.號碼() == 玩家.號碼() {
					continue
				}
				可毒的玩家號碼 = append(可毒的玩家號碼, 玩家.號碼())
			}

			uid := uuid.New().String()
			我.遊戲.旁白有話對單個玩家說(我, 傳輸資料{
				UID:     uid,
				Display: "你要毒誰呢？",
				Action:  選擇玩家,
				Data:    可毒的玩家號碼,
			}, 0)

			for {
				被女巫毒的人, err := 我.等待動作(選擇玩家, uid)
				if err != nil {
					return
				}

				玩家號碼 := 0
				err = json.Unmarshal([]byte(被女巫毒的人.Reply), &玩家號碼)
				if err != nil {
					continue
				}

				玩家, 存在 := 我.遊戲.玩家資料(玩家號碼)
				if 存在 {
					我.遊戲.殺玩家(毒殺, 玩家)
					return
				}
			}
		}
	} else {
		我.遊戲.等一下(random(3) * 2500)
	}

	return
}

func (我 *Witch) 需要夜晚行動() bool {
	return true
}
