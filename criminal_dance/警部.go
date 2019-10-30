package dance

import "strconv"

// NewPolice 新增偵探
func NewPolice(owner Player) *Police {
	return &Police{
		NewBasicCard(owner),
	}
}

// Police 偵探
type Police struct {
	*BasicCard
}

// Name 卡片名稱
func (*Police) Name() CardName {
	return 警部
}

// Detail 詳細說明
func (*Police) Detail() string {
	return `
		只有手牌在3張以下時才可以打出此牌。
		選擇一名玩家，
		當此玩家成為犯人時，你獲得遊戲勝利。
		(此卡牌在遊戲中一直處於牌面朝向其他玩家的方向)
	`
}

// CanUse 可以使用?
func (c *Police) CanUse() bool {
	return c.BasicCard.CanUse() && len(c.Owner().Cards()) <= 3
}

// Skill 功能
func (c *Police) Skill(game *Game) GameResult {

	me := c.Owner()
	data := map[string]int{}
	players := game.OtherPlayers(me)
	for _, player := range players {
		data[player.Name()] = player.No()
	}

	if len(data) == 0 {
		return 進行中
	}

	game.旁白有話對單個玩家說(me, TransferData{
		Sound:  "警長，你要預測誰是犯人?",
		Action: 等待回應,
		Data:   data,
	}, 1000)

	for {
		var no int
		td, err := me.WaitingAction(等待回應)
		if err != nil {
			no = RandomPlayerNo(players)
		} else {
			no, err = strconv.Atoi(td.Reply)
			if err != nil {
				continue
			}
		}

		he, exists := game.玩家資料(no)
		if !exists {
			continue
		}

		game.GuessFutureCriminal(me, he)
		game.旁白(TransferData{Sound: "警長預測" + he.Name() + "是犯人"}, 2000)
		return 進行中
	}
}
