package werewolf

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
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
	case 女巫:
		return NewWitch(遊戲, 號碼)
	case 狼王:
		return NewWolfKing(遊戲, 號碼)
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
		case 狼王:
			狼人玩家們 = append(狼人玩家們, 玩家.(*WolfKing))
		case 騎士:
			神職玩家們 = append(神職玩家們, 玩家.(*Knight))
		case 預言家:
			神職玩家們 = append(神職玩家們, 玩家.(*Prophesier))
		case 獵人:
			神職玩家們 = append(神職玩家們, 玩家.(*Hunter))
		case 女巫:
			神職玩家們 = append(神職玩家們, 玩家.(*Witch))
		}
	}

	return
}

// RuleOptions 角色選單
func RuleOptions() map[RULE]GROUP {
	return map[RULE]GROUP{
		平民:  人質,
		預言家: 神職,
		獵人:  神職,
		騎士:  神職,
		狼人:  狼職,
		女巫:  神職,
		狼王:  狼職,
		// 雪狼:  狼職,
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
func waitChannelBack(傳話筒 chan 傳輸資料, 對應動作 動作) (回傳資料 傳輸資料, err error) {
	if 傳話筒 == nil {
		err = errors.New("玩家已經斷線")
		return
	}

	for {
		電報 := <-傳話筒

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

		if 回傳資料.Action != 對應動作 {
			continue
		}

		break
	}

	return
}
