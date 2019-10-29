package dance

import "strconv"

// NewDog 新增神犬
func NewDog(owner Player) *Dog {
	return &Dog{
		NewBasicCard(owner),
	}
}

// Dog 神犬
type Dog struct {
	*BasicCard
}

// Name 卡片名稱
func (*Dog) Name() CardName {
	return 神犬
}

// Detail 詳細說明
func (*Dog) Detail() string {
	return `
		指定任意一名玩家，丟棄手上一張卡牌，
		如果此玩家選擇丟棄犯人卡，則神犬獲勝
		被指定的玩家有權選擇自己丟棄的卡牌，
		並在丟棄卡牌以後，
		把神犬卡收作自己的手牌。
	`
}

// CanUse 可以使用?
func (c *Dog) CanUse() bool {
	return c.BasicCard.CanUse()
}

// Skill 功能
func (c *Dog) Skill(game *Game) GameResult {

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
		Sound:  "神犬，你要咬誰?",
		Action: 等待回應,
		Data:   data,
	}, 1000)

	for {
		var no int
		td, err := me.WaitingAction(等待回應)
		if err != nil {
			no = RandomHasCardPlayerNo(players, me)
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

		game.旁白(TransferData{Sound: me.Name() + "咬了" + he.Name()}, 2000)

		for {
			// card := he.PlayCard(nil)
			card := he.TurnMe(me.No())
			if card != nil {
				game.旁白(TransferData{Sound: he.Name() + "丟棄了" + string(card.Name())}, 2000)
				if card.Name() == 犯人 {
					he.BecomeCriminal(true)
					me.BecomeDetective(true)
					return 神犬勝利
				}
				he.TakeCard(c)
				break
			}
		}

		return 進行中
	}
}
