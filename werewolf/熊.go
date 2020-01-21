package werewolf

// NewBear 建立新Bear
func NewBear(遊戲 *Game, 位子 int) *Bear {
	human := NewHuman(遊戲, 位子)
	return &Bear{
		Human: human,
	}
}

// Bear 玩家
type Bear struct {
	*Human
}

func (我 *Bear) 種族() GROUP {
	return 神職
}

func (我 *Bear) 職業() RULE {
	return 熊
}

func (我 *Bear) 需要夜晚行動() bool {
	return true
}

func (我 *Bear) 能力(i int) (_ Player) {

	if 我.出局了() {
		return
	}

	目前存活玩家們 := 我.遊戲.存活玩家們()

FIND:
	var 左邊玩家, 右邊玩家 Player
	for i := range 目前存活玩家們 {
		玩家 := 目前存活玩家們[i]
		if 我.號碼() == 玩家.號碼() {
			left := i + 1
			right := i - 1
			if i == 0 {
				right = len(目前存活玩家們) - 1
			} else if i == len(目前存活玩家們)-1 {
				left = 0
			}
			左邊玩家 = 目前存活玩家們[left]
			右邊玩家 = 目前存活玩家們[right]
			break
		}
	}

	if 左邊玩家 == nil || 右邊玩家 == nil {
		goto FIND
	}

	if (左邊玩家.種族() == 狼職 && 左邊玩家.職業() != 雪狼) ||
		(右邊玩家.種族() == 狼職 && 右邊玩家.職業() != 雪狼) {
		我.遊戲.有熊叫 = true
	} else {
		我.遊戲.有熊叫 = false
	}

	return
}
