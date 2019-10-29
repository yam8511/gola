package dance

import (
	"encoding/json"
	"errors"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

// NewCard 新增卡片
func NewCard(name CardName, owner Player) Card {
	switch name {
	case 第一發現者:
		return NewFirstFinder(owner)
	case 犯人:
		return NewCriminal(owner)
	case 不在場證明:
		return NewAbsence(owner)
	case 共犯:
		return NewAccomplice(owner)
	case 偵探:
		return NewDetective(owner)
	case 普通人:
		return NewPeople(owner)
	case 目擊者:
		return NewWitness(owner)
	case 交易:
		return NewTrade(owner)
	case 情報交換:
		return NewInformationExchange(owner)
	case 謠言:
		return NewRumor(owner)
	case 神犬:
		return NewDog(owner)
	case 警部:
	case 少年:
	}
	return nil
}

// NewPlayer 新增玩家
func NewPlayer(no int, game *Game) Player {
	return NewBasicPlayer(no, game)
}

// RandomPlayerNo 隨機曲有卡片的玩家號碼
func RandomPlayerNo(players []Player) int {

	nums := []int{}
	for _, p := range players {
		nums = append(nums, p.No())
	}
	no := nums[random(len(nums))-1]

	return no
}

// CardsInfoOutput 卡片資訊顯示
func CardsInfoOutput(cards []Card) []CardDisplay {
	myCards := []CardDisplay{}
	if len(cards) > 0 {
		for i, card := range cards {
			myCards = append(myCards, CardDisplay{
				Index:   i,
				Name:    card.Name(),
				Detail:  card.Detail(),
				Disable: !card.CanUse(),
			})
		}
	}
	return myCards
}

// CardInfoOutput 卡片資訊顯示
func CardInfoOutput(card Card) CardDisplay {
	return CardDisplay{
		Name:    card.Name(),
		Detail:  card.Detail(),
		Disable: !card.CanUse(),
	}
}

// CopyBasicCards 複製基本牌組
func CopyBasicCards() (copyCards map[CardName]int) {
	copyCards = map[CardName]int{}
	for name, count := range 基本牌組 {
		copyCards[name] = count
	}
	return
}

// PickRandomCard 隨機挑出卡片
func PickRandomCard(basicCardsMap map[CardName]int, requiredCount int, advanced bool) (cards []Card) {

	cards = []Card{}
	if len(basicCardsMap) == 0 {
		return
	}

	basicCard := []CardName{}
	for name, count := range basicCardsMap {
		for i := 0; i < count; i++ {
			basicCard = append(basicCard, name)
		}
	}

	for i := 0; i < random(999); i++ {
		rand.Shuffle(len(basicCard), func(i, j int) {
			basicCard[i], basicCard[j] = basicCard[j], basicCard[i]
		})
	}

	var currentCount int
	basicLen := len(basicCard)
	for i := 0; i < basicLen && currentCount < requiredCount; i++ {
		cardName := basicCard[i]
		if (cardName == 隨機 || cardName == 少年 || cardName == 警部) && !advanced {
			continue
		}

		card := NewCard(cardName, nil)
		if card == nil {
			continue
		}
		cards = append(cards, card)
		currentCount++
	}

	return cards
}

// Shuffle 洗牌
func Shuffle(card []Card) []Card {
	for i := 0; i < random(999); i++ {
		rand.Shuffle(len(card), func(i, j int) {
			card[i], card[j] = card[j], card[i]
		})
	}

	return card
}

// ShuffleAndCopy 洗牌並且複製一份
func ShuffleAndCopy(card []Card) []Card {

	card = Shuffle(card)

	copyCard := make([]Card, len(card))
	copy(copyCard, card)

	return copyCard
}

func random(n int) int {
	rand.Seed(time.Now().UnixNano() / 347000)
	return rand.Intn(n) + 1
}

// waitSocketBack 等待Socket回傳
func waitSocketBack(連線 *websocket.Conn, 對應動作 動作) (回傳資料 TransferData, err error) {
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
func waitChannelBack(傳話筒 chan TransferData, 對應動作 動作) (回傳資料 TransferData, err error) {
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
