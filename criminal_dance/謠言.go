package dance

import (
	"sync"
)

// NewRumor 新增謠言
func NewRumor(owner Player) *Rumor {
	return &Rumor{
		NewBasicCard(owner),
	}
}

// Rumor 謠言
type Rumor struct {
	*BasicCard
}

// Name 卡片名稱
func (*Rumor) Name() CardName {
	return 謠言
}

// Detail 詳細說明
func (*Rumor) Detail() string {
	return `
		由打出此卡牌的玩家開始，
		向右手邊玩家抽一張卡
		(先放在桌前不收入手上)，
		直到所有人都抽取一張卡牌為止
		無手牌的玩家也參加
	`
}

// CanUse 可以使用?
func (c *Rumor) CanUse() bool {
	return c.BasicCard.CanUse()
}

// Skill 卡牌技能
func (c *Rumor) Skill(game *Game) GameResult {

	game.旁白(TransferData{
		Sound:   "謠言，請向右手邊的玩家抽一張卡",
		Display: "謠言，請向右手邊的玩家抽一張卡 (右手邊 = 後一號碼)",
	}, 3000)

	drawedCards := map[int]Card{}
	takeCards := map[int]Card{}
	players := game.Players()

	wg := sync.WaitGroup{}
	mx := sync.RWMutex{}

	Draw := func(player Player) {
		defer wg.Done()

		no := player.No()
		nextNo := no + 1
		if nextNo > len(players) {
			nextNo = 1
		}

		nextPlayer, exists := game.玩家資料(nextNo)
		if !exists {
			return
		}

		card := player.PlayCard(nextPlayer)
		player.TakeCard(card)
		mx.Lock()
		drawedCards[nextNo] = card
		takeCards[no] = card
		mx.Unlock()
	}

	for _, player := range players {
		wg.Add(1)
		go Draw(player)
	}

	wg.Wait()

	ptds := map[Player]TransferData{}
	for _, player := range players {
		var takeCardOutput, drawedCardOutput *CardDisplay
		takeCard := takeCards[player.No()]
		drawedCard := drawedCards[player.No()]
		if takeCard != nil {
			t := CardInfoOutput(takeCard)
			takeCardOutput = &t
		}
		if drawedCard != nil {
			t := CardInfoOutput(drawedCard)
			drawedCardOutput = &t
		}
		ptds[player] = TransferData{
			Sound:  "請看抽到的牌",
			Action: 顯示被抽牌,
			Data: map[string]interface{}{
				"DrawedCard": drawedCardOutput,
				"GetCard":    takeCardOutput,
				"MyCard":     CardsInfoOutput(player.Cards()),
			},
		}
	}
	game.旁白有話對大家說(ptds, 3000)

	for _, player := range players {
		player.WaitingAction(顯示被抽牌)
	}

	return 進行中
}
