package dance

import (
	"sync"
)

// NewInformationExchange 新增情報交換
func NewInformationExchange(owner Player) *InformationExchange {
	return &InformationExchange{
		NewBasicCard(owner),
	}
}

// InformationExchange 情報交換
type InformationExchange struct {
	*BasicCard
}

// Name 卡片名稱
func (*InformationExchange) Name() CardName {
	return 情報交換
}

// Detail 詳細說明
func (*InformationExchange) Detail() string {
	return `
		由打出此卡牌的玩家開始，
		每位玩家給左手邊玩家傳遞一張卡牌
		(傳遞的卡牌放在桌面暫時不收入手牌)，
		直到所有玩家都抽取一張卡牌為止
		無手牌的玩家也參加
	`
}

// CanUse 可以使用?
func (c *InformationExchange) CanUse() bool {
	return c.BasicCard.CanUse()
}

// Skill 卡牌技能
func (c *InformationExchange) Skill(game *Game) GameResult {

	game.旁白(TransferData{
		Sound:   "情報交換，請選擇一張手牌給左手邊的玩家",
		Display: "情報交換，請選擇一張手牌給左手邊的玩家 (左手邊 = 前一號碼)",
		// Action:  等待回應,
	}, 3000)

	tmpCard := map[int]Card{}
	players := game.Players()

	wg := sync.WaitGroup{}
	mx := sync.RWMutex{}

	Send := func(player Player) {
		defer wg.Done()

		no := player.No()
		// 左手邊的號碼
		prevNo := no - 1
		// 如果是第一位，直接傳給最後一位
		if prevNo == 0 {
			prevNo = len(players)
		}

		card := player.PlayCard(nil, false)

		mx.Lock()
		tmpCard[prevNo] = card
		mx.Unlock()
	}

	for _, player := range players {
		wg.Add(1)
		go Send(player)
	}

	wg.Wait()

	ptds := map[Player]TransferData{}
	for _, player := range players {
		card := tmpCard[player.No()]
		var cardOutput *CardDisplay
		if card != nil {
			player.TakeCard(card)
			t := CardInfoOutput(card)
			cardOutput = &t
		}
		ptds[player] = TransferData{
			Sound:  "請看交換到的牌",
			Action: 顯示拿牌,
			Data: map[string]interface{}{
				"GetCard": cardOutput,
				"MyCard":  CardsInfoOutput(player.Cards()),
			},
		}
	}
	game.旁白有話對大家說(ptds, 3000)

	for _, player := range players {
		player.WaitingAction(顯示拿牌)
	}

	return 進行中
}
