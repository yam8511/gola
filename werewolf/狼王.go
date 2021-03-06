package werewolf

// NewKingWolf 建立新KingWolf
func NewKingWolf(遊戲 *Game, 位子 int) *KingWolf {
	wolf := NewWolf(遊戲, 位子)
	return &KingWolf{
		Wolf: wolf,
	}
}

// KingWolf 狼王
type KingWolf struct {
	*Wolf
}

func (我 *KingWolf) 職業() RULE {
	return 狼王
}

func (我 *KingWolf) 發言(投票發言 bool) bool {

	if 投票發言 {
		我.Human.發言(投票發言)
		return false
	}

	uid := newUID()
	我.遊戲.旁白有話對單個玩家說(我, 傳輸資料{
		UID:     uid,
		Display: "您要發動技能嗎? " + 我.遊戲.提示發言(),
		Action:  等待回應,
		Data: map[string]string{
			"發動✅": "yes",
			"跳過❌": "no",
		},
	}, 0)

	so, err := 我.等待動作(等待回應, uid)
	if err == nil && so.Reply == "yes" {
		return 狼人自爆(我, 我.遊戲)
	}

	return false
}

func (我 *KingWolf) 出局(殺法 KILL) {
	我.Human.出局(殺法)
	我.遊戲.判斷勝負(false)
	if 殺法 != 毒殺 && 我.遊戲.勝負 == 進行中 {
		我.能力(1)
	}
}

func (我 *KingWolf) 能力(i int) (_ Player) {
	switch i {
	case 1:
		獵人技能(我, 我.遊戲)
		return
	default:
		return 我.Wolf.能力(0)
	}
}
