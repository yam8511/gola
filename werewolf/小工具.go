package werewolf

import (
	"log"
	"math/rand"
	"time"
)

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

// PickSkiller 取出職能者
func PickSkiller(玩家們 map[string]Player) (狼人玩家們, 神職玩家們 []Skiller) {
	狼人玩家們 = []Skiller{}
	神職玩家們 = []Skiller{}

	for i := range 玩家們 {
		玩家 := 玩家們[i]
		switch 玩家.職業() {
		case 狼人:
			狼人玩家們 = append(狼人玩家們, 玩家.(*Wolf))
		case 騎士:
			神職玩家們 = append(神職玩家們, 玩家.(*Knight))
		case 預言家:
			神職玩家們 = append(神職玩家們, 玩家.(*Prophesier))
		case 獵人:
			神職玩家們 = append(神職玩家們, 玩家.(*Hunter))
		}
	}

	return
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

func 亂數洗牌(職業牌 []RULE) []RULE {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(職業牌), func(i, j int) {
		職業牌[i], 職業牌[j] = 職業牌[j], 職業牌[i]
	})
	log.Println("職業牌", 職業牌)
	return 職業牌
}
