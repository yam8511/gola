package werewolf

import (
	"encoding/json"
	"math/rand"
	"time"
)

// NewSnowWolf 建立新SnowWolf
func NewSnowWolf(遊戲 *Game, 位子 int) *SnowWolf {
	human := NewHuman(遊戲, 位子)
	return &SnowWolf{
		Human: human,
	}
}

// SnowWolf 玩家
type SnowWolf struct {
	*Human
}

func (我 *SnowWolf) 種族() GROUP {
	return 狼職
}

func (我 *SnowWolf) 職業() RULE {
	return 雪狼
}

func (我 *SnowWolf) 能力() (_ Player) {
	可殺的玩家號碼 := []int{}
	目前存活玩家們 := 我.遊戲.存活玩家們()
	for i := range 目前存活玩家們 {
		玩家 := 目前存活玩家們[i]
		if 我.遊戲.上一晚狼殺玩家號碼 == 玩家.號碼() {
			continue
		}
		可殺的玩家號碼 = append(可殺的玩家號碼, 玩家.號碼())
	}

	我.遊戲.旁白有話對單個玩家說(我, 傳輸資料{
		Display: "請問你要殺誰？",
		Action:  選擇玩家,
		Data:    可殺的玩家號碼,
	}, 0)

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

func (我 *SnowWolf) 需要夜晚行動() bool {
	return true
}
