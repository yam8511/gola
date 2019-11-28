package werewolf

// NewHunter 建立新Hunter
func NewHunter(遊戲 *Game, 位子 int) *Hunter {
	human := NewHuman(遊戲, 位子)
	return &Hunter{
		Human: human,
	}
}

// Hunter 玩家
type Hunter struct {
	*Human
}

func (我 *Hunter) 種族() GROUP {
	return 神職
}

func (我 *Hunter) 職業() RULE {
	return 獵人
}

func (我 *Hunter) 需要夜晚行動() bool {
	return false
}

func (我 *Hunter) 出局(殺法 KILL) {
	我.Human.出局(殺法)
	我.遊戲.判斷勝負(false)
	if 殺法 != 毒殺 && 我.遊戲.勝負 == 進行中 {
		我.能力()
	}
}

func (我 *Hunter) 能力() (_ Player) {
	獵人技能(我, 我.遊戲)
	return
}
