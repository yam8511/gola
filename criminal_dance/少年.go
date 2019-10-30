package dance

// NewTeenager 新增目擊者
func NewTeenager(owner Player) *Teenager {
	return &Teenager{
		NewBasicCard(owner),
	}
}

// Teenager 目擊者
type Teenager struct {
	*BasicCard
}

// Name 卡片名稱
func (*Teenager) Name() CardName {
	return 少年
}

// Detail 詳細說明
func (*Teenager) Detail() string {
	return `
		所有玩家閉上眼睛之後，
		你和「拿犯人卡」的的玩家睜開眼，再閉上。
		然後示意所有玩家睜開眼睛。
	`
}

// CanUse 可以使用?
func (c *Teenager) CanUse() bool {
	return c.BasicCard.CanUse()
}

// Skill 卡牌技能
func (c *Teenager) Skill(game *Game) GameResult {

	me := c.Owner()

	var criminalPlayer string
	var foundCriminal bool
	cards := me.Cards()
	for _, card := range cards {
		if card.Name() == 犯人 {
			foundCriminal = true
			criminalPlayer = me.Name()
			break
		}
	}

	if !foundCriminal {
		players := game.OtherPlayers(me)
		for _, player := range players {
			cards := player.Cards()
			for _, card := range cards {
				if card.Name() == 犯人 {
					foundCriminal = true
					criminalPlayer = player.Name()
					break
				}
			}

			if foundCriminal {
				break
			}
		}
	}

	game.旁白有話對單個玩家說(me, TransferData{
		Display: "少年，" + criminalPlayer + "是犯人!",
		Action:  等待回應,
	}, 1000)

	me.WaitingAction(等待回應)
	return 進行中
}
