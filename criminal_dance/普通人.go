package dance

// NewPeople 新增普通人
func NewPeople(owner Player) *People {
	return &People{
		NewBasicCard(owner),
	}
}

// People 普通人
type People struct {
	*BasicCard
}

// Name 卡片名稱
func (*People) Name() CardName {
	return 普通人
}

// Detail 詳細說明
func (*People) Detail() string {
	return `打出時沒有任何事情發生`
}

// CanUse 可以使用?
func (c *People) CanUse() bool {
	return c.BasicCard.CanUse()
}

// Skill 卡牌技能
func (*People) Skill(game *Game) GameResult {
	return 進行中
}
