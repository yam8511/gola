package werewolf

import (
	"encoding/json"
	"strconv"
)

func 獵人技能(我 Skiller, 遊戲 *Game) {

	遊戲.旁白(傳輸資料{Sound: strconv.Itoa(我.號碼()) + "號玩家啟動角色技能，請問你要帶走誰？"}, 3000)

	可殺的玩家號碼 := map[string]int{"不帶": -1}
	目前存活玩家們 := 遊戲.存活玩家們()
	for i := range 目前存活玩家們 {
		號碼 := 目前存活玩家們[i].號碼()
		可殺的玩家號碼[strconv.Itoa(號碼)] = 號碼
	}

	uid := newUID()
	遊戲.旁白有話對單個玩家說(我, 傳輸資料{
		UID:     uid,
		Display: "請問你要帶走誰？",
		Action:  等待回應,
		Data:    可殺的玩家號碼,
	}, 1000)

	for {
		so, err := 我.等待動作(等待回應, uid)
		if err != nil {
			return
		}

		玩家號碼 := 0
		err = json.Unmarshal([]byte(so.Reply), &玩家號碼)
		if err != nil {
			continue
		}

		if 玩家號碼 == -1 {
			return
		}

		玩家, 存在 := 遊戲.玩家資料(玩家號碼)
		if 存在 {
			uid := newUID()
			遊戲.旁白(傳輸資料{
				UID:    uid,
				Sound:  strconv.Itoa(玩家.號碼()) + "號玩家淘汰!" + 遊戲.提示點擊(false),
				Action: 等待回應,
			}, 3000)

			for _, 存活玩家 := range 目前存活玩家們 {
				存活玩家.等待動作(等待回應, uid)
			}

			遊戲.殺玩家(獵殺, 玩家)
			遊戲.玩家出局()

			return
		}
	}
}

func 狼人自爆(我 Skiller, 遊戲 *Game) bool {
	if !遊戲.可自爆 {
		return false
	}

	遊戲.時間 = 投票

	uid := newUID()
	遊戲.旁白(傳輸資料{
		UID:    uid,
		Sound:  strconv.Itoa(我.號碼()) + "號玩家自爆。" + 遊戲.提示點擊(true),
		Action: 等待回應,
	}, 3000)

	for _, 存活玩家 := range 遊戲.存活玩家們() {
		存活玩家.等待動作(等待回應, uid)
	}

	遊戲.殺玩家(自爆, 我)
	遊戲.玩家出局()

	return true
}
