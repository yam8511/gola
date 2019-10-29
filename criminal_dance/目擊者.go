package dance

import "strconv"

// NewWitness 新增目擊者
func NewWitness(owner Player) *Witness {
	return &Witness{
		NewBasicCard(owner),
	}
}

// Witness 目擊者
type Witness struct {
	*BasicCard
}

// Name 卡片名稱
func (*Witness) Name() CardName {
	return 目擊者
}

// Detail 詳細說明
func (*Witness) Detail() string {
	return `
		指定任意一名玩家
		看該玩家手上的所有手牌
	`
}

// CanUse 可以使用?
func (c *Witness) CanUse() bool {
	return c.BasicCard.CanUse()
}

// Skill 卡牌技能
func (c *Witness) Skill(game *Game) GameResult {

	me := c.Owner()
	data := map[string]int{}
	players := game.HasOtherCardPlayers(me)
	for _, player := range players {
		data[player.Name()] = player.No()
	}

	if len(data) == 0 {
		return 進行中
	}

	game.旁白有話對單個玩家說(me, TransferData{
		Sound:  "目擊者，你要看誰的手牌?",
		Action: 等待回應,
		Data:   data,
	}, 1000)

	for {
		td, err := me.WaitingAction(等待回應)
		if err != nil {
			return 進行中
		}

		no, err := strconv.Atoi(td.Reply)
		if err != nil {
			continue
		}

		he, exists := game.玩家資料(no)
		if !exists {
			continue
		}

		game.旁白(TransferData{Sound: me.Name() + "目擊了" + he.Name()}, 2000)

		cards := []CardName{}
		for _, card := range he.Cards() {
			cards = append(cards, card.Name())
		}

		game.旁白有話對單個玩家說(me, TransferData{
			Sound:  he.Name() + "的手牌",
			Action: 看手牌,
			Data:   cards,
		}, 100)

		me.WaitingAction(看手牌)

		return 進行中
	}
}
