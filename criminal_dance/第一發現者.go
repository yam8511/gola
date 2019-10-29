package dance

// NewFirstFinder 新增第一發現者
func NewFirstFinder(owner Player) *FirstFinder {
	return &FirstFinder{
		NewBasicCard(owner),
	}
}

// FirstFinder 第一發現者
type FirstFinder struct {
	*BasicCard
}

// Name 卡片名稱
func (*FirstFinder) Name() CardName {
	return 第一發現者
}

// Detail 詳細說明
func (*FirstFinder) Detail() string {
	return `
		拿到此牌的人，
		就是第一張打出來的起始玩家。
	`
}

// CanUse 可以使用?
func (c *FirstFinder) CanUse() bool {
	return true
}

// Skill 卡牌技能
func (*FirstFinder) Skill(game *Game) GameResult {
	return 進行中
}
