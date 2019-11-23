package werewolf

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// NewPlayer 產生新角色
func NewPlayer(角色 RULE, 遊戲 *Game, 號碼 int) Player {
	switch 角色 {
	case 平民:
		return NewHuman(遊戲, 號碼)
	case 騎士:
		return NewKnight(遊戲, 號碼)
	case 預言家:
		return NewSeer(遊戲, 號碼)
	case 獵人:
		return NewHunter(遊戲, 號碼)
	case 女巫:
		return NewWitch(遊戲, 號碼)
	case 狼人:
		return NewWolf(遊戲, 號碼)
	case 狼王:
		return NewWolfKing(遊戲, 號碼)
	case 雪狼:
		return NewSnowWolf(遊戲, 號碼)
	}

	return nil
}

func 角色選單() map[RULE]GROUP {
	return map[RULE]GROUP{
		平民:  人質,
		預言家: 神職,
		獵人:  神職,
		騎士:  神職,
		狼人:  狼職,
		女巫:  神職,
		狼王:  狼職,
		雪狼:  狼職,
	}
}

func 快速組合() map[string]map[RULE]int {
	return map[string]map[RULE]int{
		"4": map[RULE]int{
			平民:  2,
			狼人:  1,
			預言家: 1,
		},
		"5": map[RULE]int{
			平民:  2,
			狼人:  1,
			預言家: 1,
			女巫:  1,
		},
		"6": map[RULE]int{
			平民:  2,
			狼人:  2,
			預言家: 1,
			女巫:  1,
		},
		"7": map[RULE]int{
			平民:  2,
			狼人:  2,
			預言家: 1,
			女巫:  1,
			獵人:  1,
		},
		"8": map[RULE]int{
			平民:  3,
			狼人:  2,
			預言家: 1,
			女巫:  1,
			獵人:  1,
		},
		"9": map[RULE]int{
			平民:  3,
			狼人:  3,
			預言家: 1,
			女巫:  1,
			獵人:  1,
		},
		"10(狼王)": map[RULE]int{
			平民:  3,
			狼人:  2,
			狼王:  1,
			預言家: 1,
			女巫:  1,
			獵人:  1,
			騎士:  1,
		},
		"10(雪狼)": map[RULE]int{
			平民:  3,
			狼人:  2,
			雪狼:  1,
			預言家: 1,
			女巫:  1,
			獵人:  1,
			騎士:  1,
		},
	}
}

// PickSkiller 取出職能者
func PickSkiller(玩家們 map[string]Player) (狼人玩家們, 神職玩家們 []Skiller) {
	狼人玩家們 = []Skiller{}
	神職玩家們 = []Skiller{}

	for i := range 玩家們 {
		玩家 := 玩家們[i]
		switch 玩家.種族() {
		case 狼職:
			狼人玩家們 = append(狼人玩家們, 玩家.(Skiller))
		case 神職:
			神職玩家們 = append(神職玩家們, 玩家.(Skiller))
		}
	}

	return
}

func 亂數洗牌(職業牌 []RULE) []RULE {
	for i := 0; i < random(999); i++ {
		rand.Shuffle(len(職業牌), func(i, j int) {
			職業牌[i], 職業牌[j] = 職業牌[j], 職業牌[i]
		})
	}
	log.Println("職業牌", 職業牌)
	return 職業牌
}

func random(n int) int {
	rand.Seed(time.Now().UnixNano() / 347000)
	return rand.Intn(n) + 1
}

func newUID() string {
	return uuid.New().String()
}

// waitSocketBack 等待Socket回傳
func waitSocketBack(連線 *websocket.Conn, 對應動作 動作) (回傳資料 傳輸資料, err error) {
	var msg []byte

	for {
		_, msg, err = 連線.ReadMessage()
		if err != nil {
			return
		}

		// 如果沒有對應動作，直接把資料傳給傳話筒用
		if 對應動作 == 無 {
			回傳資料.Action = 給傳話筒
			回傳資料.Reply = string(msg)
			return
		}

		err = json.Unmarshal(msg, &回傳資料)
		if err != nil {
			continue
		}

		if 回傳資料.Action != 對應動作 {
			continue
		}

		break
	}

	return
}

// waitChannelBack 等待Channel回傳
func waitChannelBack(傳話筒 chan 傳輸資料, 對應動作 動作, uid string) (回傳資料 傳輸資料, err error) {
	if 傳話筒 == nil {
		err = errors.New("玩家已經斷線")
		return
	}

	for {
		var 電報 傳輸資料
		var ok bool
		select {
		case 電報, ok = <-傳話筒:
			if !ok {
				err = errors.New("玩家已經斷線")
				return
			}
		case <-time.After(time.Second):
			runtime.Gosched()
			continue
		}

		if 電報.Action == 無 {
			err = errors.New("玩家已經斷線")
			return
		}

		if 對應動作 == 無 {
			回傳資料 = 電報
			return
		}

		if 電報.Action != 給傳話筒 {
			continue
		}

		err = json.Unmarshal([]byte(電報.Reply), &回傳資料)
		if err != nil {
			continue
		}

		if 回傳資料.Action != 對應動作 || 回傳資料.UID != uid {
			continue
		}

		break
	}

	return
}
