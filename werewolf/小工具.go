package werewolf

import (
	"math/rand"
)

func 亂數洗牌(array []RULE) []RULE {
	for i := len(array) - 1; i >= 0; i-- {
		p := 區間隨機數(0, int64(i))
		a := array[i]
		array[i] = array[p]
		array[p] = a
	}
	return array
}

// 區間隨機數 區間隨機數
func 區間隨機數(min, max int64) int64 {
	if min >= max || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}
