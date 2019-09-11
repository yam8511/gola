package werewolf

import (
	"log"
	"math/rand"
	"time"
)

func 亂數洗牌(職業牌 []RULE) []RULE {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(職業牌), func(i, j int) {
		職業牌[i], 職業牌[j] = 職業牌[j], 職業牌[i]
	})
	log.Println("職業牌", 職業牌)
	return 職業牌
}
