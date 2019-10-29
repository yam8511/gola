package dance

import "strconv"

// NewDetective 新增偵探
func NewDetective(owner Player) *Detective {
	return &Detective{
		NewBasicCard(owner),
	}
}

// Detective 偵探
type Detective struct {
	*BasicCard
}

// Name 卡片名稱
func (*Detective) Name() CardName {
	return 偵探
}

// Detail 詳細說明
func (*Detective) Detail() string {
	return `
		指定任意一名玩家，
		猜他手牌中有犯人卡
		偵探牌只能在手牌3張以下時，
		才可以打出。
		猜中犯人時贏得遊戲。
	`
}

// CanUse 可以使用?
func (c *Detective) CanUse() bool {
	return c.BasicCard.CanUse() && len(c.Owner().Cards()) <= 3
}

// Skill 功能
func (c *Detective) Skill(game *Game) GameResult {

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
		Sound:  "偵探，你要指認誰?",
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

		game.旁白(TransferData{Sound: "偵探指認" + he.Name()}, 2000)

		yes := he.BeAskedCriminal()
		if yes {
			game.旁白(TransferData{Sound: he.Name() + "是犯人！"}, 2000)
			me.BecomeDetective(true)
			return 偵探勝利
		}

		game.旁白(TransferData{Sound: he.Name() + "不是犯人！"}, 2000)
		return 進行中
	}
}
