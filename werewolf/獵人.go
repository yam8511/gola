package werewolf

import (
	"encoding/json"
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

func (我 *Hunter) 需要夜晚行動() bool {
	return false
}

func (我 *Hunter) 出局(killed KILL) {
	if killed != 毒殺 {
		我.能力()
	}
}

func (我 *Hunter) 能力() (_ Player) {

	我.遊戲.旁白(傳輸資料{
		Sound: "啟動角色技能，請問你要帶走誰？",
	})

	可殺的玩家號碼 := []int{}
	目前存活玩家們 := 我.遊戲.存活玩家們()
	for i := range 目前存活玩家們 {
		玩家 := 目前存活玩家們[i]
		可殺的玩家號碼 = append(可殺的玩家號碼, 玩家.號碼())
	}

	我.遊戲.旁白有話對單個玩家說(我, 傳輸資料{
		Display: "請問你要帶走誰？",
		Action:  選擇玩家,
		Data:    可殺的玩家號碼,
	})

	for {

		so, err := waitChannelBack(我.傳話筒, 選擇玩家)
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
			我.遊戲.殺玩家(獵殺, 玩家)
			return
		}
	}
}
