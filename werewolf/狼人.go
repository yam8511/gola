package werewolf

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"
)

// NewWolf 建立新Wolf
func NewWolf(位子 int) *Wolf {
	return &Wolf{
		Human: Human{
			位子: 位子,
		},
	}
}

// Wolf 玩家
type Wolf struct {
	Human
}

func (我 *Wolf) 種族() GROUP {
	return 狼職
}

func (我 *Wolf) 職業() RULE {
	return 狼人
}

func (我 *Wolf) 能力() {
	可殺的玩家號碼 := []int{}
	for i := range 我.遊戲.存活玩家們() {
		可殺的玩家號碼 = append(可殺的玩家號碼, 我.遊戲.玩家們[i].號碼())
	}

	if 我.連線 == nil {
		rand.Seed(time.Now().UTC().UnixNano())
		rd := rand.Int()
		玩家號碼 := 可殺的玩家號碼[rd%len(可殺的玩家號碼)]
		for i := range 我.遊戲.玩家們 {
			玩家 := 我.遊戲.玩家們[i]
			if 玩家.號碼() == 玩家號碼 {
				我.遊戲.殺玩家(狼殺, 玩家)
				return
			}
		}
		return
	}

	for {
		err := 我.連線.WriteJSON(map[string]interface{}{
			"event": "請問你要殺誰？",
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
					我.遊戲.殺玩家(狼殺, 玩家)
					return
				}
			}
		}
	}
}
