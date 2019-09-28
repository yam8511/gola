package werewolf

import (
	"encoding/json"
	"math/rand"
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

func (我 *Wolf) 能力() (_ Player) {
	可殺的玩家號碼 := []int{}
	目前存活玩家們 := 我.遊戲.存活玩家們()
	for i := range 目前存活玩家們 {
		玩家 := 目前存活玩家們[i]
		可殺的玩家號碼 = append(可殺的玩家號碼, 玩家.號碼())
	}

	我.遊戲.旁白(傳輸資料{
		Sound:  "請問你要殺誰？",
		Action: 選擇玩家,
		Data:   可殺的玩家號碼,
	})

	for {

		so, err := waitChannelBack(我.傳話筒, 選擇玩家)
		if err != nil {
			rand.Seed(time.Now().UTC().UnixNano())
			rd := rand.Int()
			玩家號碼 := 可殺的玩家號碼[rd%len(可殺的玩家號碼)]
			玩家, 存在 := 我.遊戲.玩家資料(玩家號碼)
			if 存在 {
				// 我.遊戲.殺玩家(狼殺, 玩家)
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
