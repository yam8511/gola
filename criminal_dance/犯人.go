package dance

// NewCriminal 新增第一發現者
func NewCriminal(owner Player) *Criminal {
	return &Criminal{
		NewBasicCard(owner),
	}
}

// Criminal 第一發現者
type Criminal struct {
	*BasicCard
}

// Name 卡片名稱
func (*Criminal) Name() CardName {
	return 犯人
}

// Detail 詳細說明
func (*Criminal) Detail() string {
	return `
		拿到犯人的玩家，
		只能在手上僅有此牌時
		才可以打出此牌，
		並同時獲得勝利。
	`
}

// CanUse 可以使用?
func (c *Criminal) CanUse() bool {
	return len(c.Owner().Cards()) == 1
}

// Skill 功能
func (c *Criminal) Skill(game *Game) GameResult {
	c.Owner().BecomeCriminal(true)
	return 犯人勝利
}
