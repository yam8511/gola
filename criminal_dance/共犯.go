package dance

// NewAccomplice 新增共犯
func NewAccomplice(owner Player) *Accomplice {
	return &Accomplice{
		NewBasicCard(owner),
	}
}

// Accomplice 共犯
type Accomplice struct {
	*BasicCard
}

// Name 卡片名稱
func (*Accomplice) Name() CardName {
	return 共犯
}

// Detail 詳細說明
func (*Accomplice) Detail() string {
	return `
		打出此卡後你就是共犯
		(在手上沒有任何效果)
		犯人勝利時，共犯也獲得勝利
		犯人失敗時，共犯也失敗。
	`
}

// CanUse 可以使用?
func (c *Accomplice) CanUse() bool {
	return c.BasicCard.CanUse()
}

// Skill 功能
func (c *Accomplice) Skill(game *Game) GameResult {
	c.Owner().BecomeAccomplice(true)
	return 進行中
}
