package werewolf

// NewPlayer 產生新角色
func NewPlayer(角色 RULE, 遊戲 *Game, 號碼 int) Player {
	switch 角色 {
	case 平民:
		return NewHuman(遊戲, 號碼)
	case 狼人:
		return NewWolf(遊戲, 號碼)
	case 騎士:
		return NewKnight(遊戲, 號碼)
	case 預言家:
		return NewProphesier(遊戲, 號碼)
	case 獵人:
		return NewHunter(遊戲, 號碼)
	}

	return nil
}

// RuleOptions 角色選單
func RuleOptions() map[RULE]GROUP {
	return map[RULE]GROUP{
		平民:  人質,
		預言家: 神職,
		女巫:  神職,
		獵人:  神職,
		騎士:  神職,
		狼人:  狼職,
		狼王:  狼職,
		雪狼:  狼職,
	}
}
