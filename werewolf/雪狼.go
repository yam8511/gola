package werewolf

// NewSnowWolf 建立新SnowWolf
func NewSnowWolf(遊戲 *Game, 位子 int) *SnowWolf {
	wolf := NewWolf(遊戲, 位子)
	return &SnowWolf{
		Wolf: wolf,
	}
}

// SnowWolf 玩家
type SnowWolf struct {
	*Wolf
}

func (我 *SnowWolf) 職業() RULE {
	return 雪狼
}
