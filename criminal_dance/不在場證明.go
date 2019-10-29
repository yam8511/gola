package dance

// NewAbsence 新增不在場證明
func NewAbsence(owner Player) *Absence {
	return &Absence{
		NewBasicCard(owner),
	}
}

// Absence 不在場證明
type Absence struct {
	*BasicCard
}

// Name 卡片名稱
func (*Absence) Name() CardName {
	return 不在場證明
}

// Detail 詳細說明
func (*Absence) Detail() string {
	return `
		擁有不再證明的人，
		被抓到時可以不必承認自己是犯人。
	`
}

// CanUse 可以使用?
func (c *Absence) CanUse() bool {
	return c.BasicCard.CanUse()
}

// Skill 卡牌技能
func (*Absence) Skill(game *Game) GameResult {
	return 進行中
}
