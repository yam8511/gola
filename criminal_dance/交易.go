package dance

import (
	"strconv"
	"sync"
)

// NewTrade 新增交易
func NewTrade(owner Player) *Trade {
	return &Trade{
		NewBasicCard(owner),
	}
}

// Trade 交易
type Trade struct {
	*BasicCard
}

// Name 卡片名稱
func (*Trade) Name() CardName {
	return 交易
}

// Detail 詳細說明
func (*Trade) Detail() string {
	return `
		指定任意一名玩家
		交換手上一張卡牌
		如果該玩家手中沒有卡牌
		則不能指定此玩家
	`
}

// CanUse 可以使用?
func (c *Trade) CanUse() bool {
	return c.BasicCard.CanUse()
}

// Skill 卡牌技能
func (c *Trade) Skill(game *Game) GameResult {

	me := c.Owner()
	data := map[string]int{}
	players := game.HasOtherCardPlayers(me)
	for _, player := range players {
		data[player.Name()] = player.No()
	}

	if len(data) == 0 {
		return 進行中
	}

	game.旁白有話對單個玩家說(me, TransferData{
		Sound:  "你要和誰交易?",
		Action: 等待回應,
		Data:   data,
	}, 1000)

	for {
		var no int
		td, err := me.WaitingAction(等待回應)
		if err != nil {
			no = RandomHasCardPlayerNo(players, me)
		} else {
			no, err = strconv.Atoi(td.Reply)
			if err != nil {
				continue
			}
		}

		he, exists := game.玩家資料(no)
		if !exists {
			continue
		}

		game.旁白(TransferData{Sound: me.Name() + "和" + he.Name() + "交易"}, 2000)

		var herCard, myCard Card

		wg := sync.WaitGroup{}
		wg.Add(2)

		go func() {
			// 對方出卡
			herCard = he.PlayCard(nil)
			wg.Done()
		}()

		go func() {
			// 換我出卡
			myCard = me.PlayCard(nil)
			wg.Done()
		}()

		wg.Wait()

		// 開始交換
		me.TakeCard(herCard)
		he.TakeCard(myCard)

		ptds := map[Player]TransferData{
			me: TransferData{
				Sound:  "請看交換到的牌",
				Action: 顯示拿牌,
				Data: map[string]interface{}{
					"GetCard": CardInfoOutput(herCard),
					"MyCard":  CardsInfoOutput(me.Cards()),
				},
			},
			he: TransferData{
				Sound:  "請看交換到的牌",
				Action: 顯示拿牌,
				Data: map[string]interface{}{
					"GetCard": CardInfoOutput(myCard),
					"MyCard":  CardsInfoOutput(he.Cards()),
				},
			},
		}
		game.旁白有話對大家說(ptds, 3000)

		he.WaitingAction(顯示拿牌)
		me.WaitingAction(顯示拿牌)

		return 進行中
	}
}
