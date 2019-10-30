package dance

import (
	"encoding/json"
	"errors"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Game 犯人在跳舞遊戲
type Game struct {
	roomMaster     Player
	players        []Player
	police         Player
	futureCriminal Player
	mx             sync.RWMutex
	cards          []Card
	targetPoint    int
	gameState      遊戲階段
}

func (g *Game) 加入(conn *websocket.Conn) {
	進入遊戲 := func(me Player) {
		g.旁白有話對單個玩家說(me, TransferData{
			Sound:  "你的角色",
			Action: 拿到角色,
			Data: map[string]interface{}{
				"No":   me.No(),
				"Name": me.Name(),
			}}, 0,
		)

		g.加入玩家(me)

		me.Waiting()
	}

	for g.目前階段() > 準備階段 {
		players := g.Players()
		availablePlayer := []int{}
		for _, player := range players {
			if !player.IsConnected() {
				availablePlayer = append(availablePlayer, player.No())
			}
		}

		err := g.旁白有話對連線說(conn, TransferData{
			Sound:  "遊戲已經開始，想重新進入遊戲，請輸入玩家號碼",
			Action: 遊戲已開始,
			Data:   availablePlayer,
		}, 0)

		if err != nil {
			return
		}

		so, err := waitSocketBack(conn, 遊戲已開始)
		if err != nil {
			return
		}

		var input struct {
			Name string `json:"name"`
			No   int    `json:"no"`
		}

		err = json.Unmarshal([]byte(so.Reply), &input)
		if err != nil {
			continue
		}

		me, exist := g.玩家資料(input.No)
		if exist && me.Join(conn, input.Name) {
			g.RefreshDesktop()
			進入遊戲(me)
			return
		}
	}

	var me Player

	for {
		pos := g.顯示可選位子()
		if len(pos) == 0 {
			g.旁白有話對連線說(conn, TransferData{Sound: "遊戲人數已額滿"}, 0)
			return
		}

		err := g.旁白有話對連線說(conn, TransferData{
			Sound:  "請選擇號碼",
			Action: 選擇號碼,
			Data:   pos,
		}, 0)

		if err != nil {
			g.mx.RLock()
			roomMaster := g.roomMaster
			g.mx.RUnlock()
			if roomMaster == nil {
				g.重置()
			}
			return
		}

		so, err := waitSocketBack(conn, 選擇號碼)
		if err != nil {
			g.mx.RLock()
			roomMaster := g.roomMaster
			g.mx.RUnlock()
			if roomMaster == nil {
				g.重置()
			}
			return
		}

		var input struct {
			Name string `json:"name"`
			No   int    `json:"no"`
		}

		err = json.Unmarshal([]byte(so.Reply), &input)
		if err != nil {
			continue
		}

		var exist bool
		me, exist = g.玩家資料(input.No)
		if exist && !me.IsConnected() && me.Join(conn, input.Name) {
			進入遊戲(me)
			return
		}
	}
}

func (g *Game) 開始() {

	reachPointPlayer := map[int]Player{}
	players := g.Players()
	playerNum := len(players)
	for {
		if g.目前階段() == 初始階段 {
			return
		}

		g.ChangeGameState(開始階段)
		// 洗牌
		cards := ShuffleAndCopy(g.cards)

		var firstPlayer Player
		// 發牌
		for i, card := range cards {
			player := players[i%playerNum]
			player.TakeCard(card)
			if player.HasFirstFinder() {
				firstPlayer = player
			}
		}
		g.RefreshDesktop()
		g.旁白(TransferData{Sound: "發完牌囉，看一下手牌"}, 0)
		g.等一下(3000)

		var result GameResult
		// 輪流出牌
		for cardPlayers := g.HasCardPlayers(); len(cardPlayers) > 0; cardPlayers = g.HasCardPlayers() {
			for _, p := range cardPlayers {
				if g.目前階段() == 初始階段 {
					return
				}

				if firstPlayer != nil {
					if p.No() != firstPlayer.No() {
						continue
					}
					firstPlayer = nil
				}
				g.旁白(TransferData{
					Sound:  p.Name() + "出牌",
					Action: 輪到我,
					Data:   p.No(),
				}, 1000)
				card := p.TurnMe(p.No())
				g.RefreshDesktop()
				if card == nil {
					continue
				}
				g.旁白(TransferData{
					Sound:  p.Name() + "出" + string(card.Name()),
					Action: 顯示出牌,
					Data: map[string]interface{}{
						"Player": p.Name(),
						"Card":   CardInfoOutput(card),
					},
				}, 3000)
				if card != nil {
					result = card.Skill(g)
					g.RefreshDesktop()
					if result != 進行中 {
						goto SETTLE
					}
				}
			}
		}

	SETTLE:

		if g.目前階段() == 初始階段 {
			return
		}

		if criminal := g.FutureCriminal(); criminal != nil {
			for _, player := range g.Players() {
				if player.No() == criminal.No() {
					if player.IsCriminal() {
						result = 警部勝利
					}
					break
				}
			}
		}

		if police := g.Police(); result == 警部勝利 && police != nil {
			for _, player := range g.Players() {
				player.BecomeDetective(false)
			}
			police.BecomeDetective(true)
		}

		// 結算分數
		g.ChangeGameState(結算階段)
		var normalPoint, criminalPoint, detectivePoint int
		switch result {
		case 偵探勝利:
			normalPoint = 1
			detectivePoint = 2
		case 神犬勝利, 警部勝利:
			normalPoint = 1
			detectivePoint = 3
		case 犯人勝利:
			criminalPoint = 2
		}

		// 結算每位玩家的分數
		playerName := map[int]string{}
		pointChanged := map[int]int{}
		playerPoint := map[int]int{}
		for _, p := range players {
			no := p.No()
			switch {
			case p.IsDetective():
				p.TakePoint(detectivePoint)
				pointChanged[no] = detectivePoint
			case p.IsCriminal() || p.IsAccomplice():
				p.TakePoint(criminalPoint)
				pointChanged[no] = criminalPoint
			default:
				p.TakePoint(normalPoint)
				pointChanged[no] = normalPoint
			}

			playerName[no] = p.Name()
			cp := p.CurrentPoint()
			playerPoint[no] = cp

			if cp >= g.targetPoint {
				reachPointPlayer[no] = p
			}
		}

		// 通知結算顯示
		g.旁白(TransferData{
			Sound:  string(result) + "，結算分數",
			Action: 顯示結算,
			Data: map[string]interface{}{
				"GameResult":   result,
				"PlayerName":   playerName,
				"PointChanged": pointChanged,
				"PlayerPoint":  playerPoint,
			},
		}, 2000)

		for _, p := range players {
			p.WaitingAction(顯示結算)
			p.ClearCard()
			p.BecomeCriminal(false)
			p.BecomeAccomplice(false)
			p.BecomeDetective(false)
		}
		g.GuessFutureCriminal(nil, nil)

		// 如果有玩家達到目標分數，遊戲結束
		if len(reachPointPlayer) > 0 {
			break
		}

		if g.目前階段() == 初始階段 {
			return
		}

		// 通知結算顯示
		g.旁白(TransferData{Sound: "開始新的一局"}, 2000)
	}

	// 遊戲結束
	if g.目前階段() == 初始階段 {
		return
	}

	winner := []string{}
	for _, player := range reachPointPlayer {
		winner = append(winner, player.Name())
	}
	sound := "遊戲結束，" + strings.Join(winner, ", ") + "獲勝"

	g.旁白(TransferData{
		Sound:  sound,
		Action: 遊戲結束,
	}, 1200)
	runtime.Gosched()

	g.重置()
}

// HasSetup 已經初始設定過
func (g *Game) HasSetup() bool {
	g.mx.Lock()
	isSetup := len(g.cards) > 0 && len(g.players) > 0
	g.mx.Unlock()
	return isSetup
}

// Setup 初始設定
func (g *Game) Setup(
	setupData 規則設定,
) {
	if len(setupData.CardSet) == 0 {
		return
	}

	g.mx.Lock()

	basicCards := CopyBasicCards()
	cards := []Card{}
	extraCards := 0
	for name, count := range setupData.CardSet {
		if count == 0 {
			continue
		}

		if name == 隨機 {
			extraCards += count
			continue
		}

		if count > basicCards[name] {
			extraCards += count - basicCards[name]
			count = basicCards[name]
		}

		for i := 0; i < count; i++ {
			card := NewCard(name, nil)
			if card == nil {
				extraCards++
				continue
			}
			cards = append(cards, card)
		}
		basicCards[name] -= count
	}

	if extraCards > 0 {
		cards = append(cards, PickRandomCard(basicCards, extraCards, setupData.Advanced)...)
	}

	g.cards = Shuffle(cards)

	players := []Player{}
	for i := 0; i < setupData.PlayerCount; i++ {
		players = append(players, NewPlayer(i+1, g))
	}
	g.players = players

	g.targetPoint = setupData.TargetPoint

	g.mx.Unlock()

	g.ChangeGameState(準備階段)
}

// Players 玩家們
func (g *Game) Players() []Player {
	g.mx.RLock()
	players := g.players
	g.mx.RUnlock()
	return players
}

// HasCardPlayers 有卡牌的玩家們
func (g *Game) HasCardPlayers() []Player {
	players := []Player{}
	g.mx.RLock()
	for _, p := range g.players {
		if !p.IsEmptyCard() {
			players = append(players, p)
		}
	}
	g.mx.RUnlock()
	return players
}

// HasOtherCardPlayers 有卡牌的玩家們
func (g *Game) HasOtherCardPlayers(onwer Player) []Player {
	players := []Player{}
	g.mx.RLock()
	for _, p := range g.players {
		if !p.IsEmptyCard() && p.No() != onwer.No() {
			players = append(players, p)
		}
	}
	g.mx.RUnlock()
	return players
}

// OtherPlayers 其他玩家們
func (g *Game) OtherPlayers(onwer Player) []Player {
	players := []Player{}
	g.mx.RLock()
	for _, p := range g.players {
		if p.No() != onwer.No() {
			players = append(players, p)
		}
	}
	g.mx.RUnlock()
	return players
}

// RefreshDesktop 刷新桌面
func (g *Game) RefreshDesktop() {
	cardNum := map[int]int{}
	throwCard := map[int][]CardDisplay{}
	playerName := map[int]string{}
	playerPoint := map[int]int{}
	playerNo := []int{}

	players := g.Players()
	for _, player := range players {
		no := player.No()
		playerNo = append(playerNo, no)
		playerName[no] = player.Name()
		playerPoint[no] = player.CurrentPoint()
		cardNum[no] = len(player.Cards())
		throwCard[no] = CardsInfoOutput(player.ThrowCards())
	}

	ptds := map[Player]TransferData{}
	for _, player := range players {
		ptds[player] = TransferData{
			Action: 更新桌面,
			Data: map[string]interface{}{
				"TargetPoint": g.targetPoint,
				"PlayerNo":    playerNo,
				"PlayerName":  playerName,
				"PlayerPoint": playerPoint,
				"CardNum":     cardNum,
				"ThrowCard":   throwCard,
				"MyCard":      CardsInfoOutput(player.Cards()),
			},
		}
	}
	g.旁白有話對大家說(ptds, 0)
}

// ChangeGameState 變更遊戲階段
func (g *Game) ChangeGameState(state 遊戲階段) {
	g.mx.Lock()
	g.gameState = state
	g.mx.Unlock()
}

// GuessFutureCriminal 預測未來的犯人
func (g *Game) GuessFutureCriminal(police, criminal Player) {
	g.mx.Lock()
	g.police = police
	g.futureCriminal = criminal
	g.mx.Unlock()
}

// FutureCriminal 預測未來的犯人
func (g *Game) FutureCriminal() (criminal Player) {
	g.mx.RLock()
	criminal = g.futureCriminal
	g.mx.RUnlock()
	return
}

// Police 警長
func (g *Game) Police() (plice Player) {
	g.mx.RLock()
	plice = g.police
	g.mx.RUnlock()
	return
}

func (g *Game) 目前階段() 遊戲階段 {
	g.mx.RLock()
	state := g.gameState
	g.mx.RUnlock()
	return state
}

func (g *Game) 玩家資料(no int) (me Player, exist bool) {
	g.mx.RLock()
	defer g.mx.RUnlock()
	for id := range g.players {
		me = g.players[id]
		if me.No() == no {
			exist = true
			break
		}
	}

	return
}

func (g *Game) 重置() {
	for _, player := range g.Players() {
		if player.IsConnected() {
			player.Conn().Close()
		}
	}

	g.mx.Lock()
	g.roomMaster = nil
	g.cards = nil
	g.players = nil
	g.gameState = 初始階段
	g.targetPoint = 0
	g.mx.Unlock()
}

func (g *Game) 旁白(台詞 TransferData, 豪秒數 int) {
	players := g.Players()
	for i := range players {
		p := players[i]
		p.TransferToMe(台詞)
	}

	g.等一下(豪秒數)
}

func (g *Game) 旁白有話對大家說(ptds map[Player]TransferData, 豪秒數 int) {
	for player, td := range ptds {
		player.TransferToMe(td)
	}
	g.等一下(豪秒數)
}

func (g *Game) 旁白有話對單個玩家說(玩家 Player, 台詞 TransferData, 豪秒數 int) {
	玩家.TransferToMe(台詞)
	g.等一下(豪秒數)
}

func (g *Game) 旁白有話對連線說(連線 *websocket.Conn, 台詞 TransferData, 豪秒數 int) error {
	if 連線 != nil {
		err := 連線.WriteJSON(台詞)
		if err != nil {
			return errors.New("WebSocket連線斷線了")
		}
	}

	g.等一下(豪秒數)
	return nil
}

func (g *Game) 等一下(豪秒數 int) {
	if 豪秒數 > 0 {
		time.Sleep(time.Millisecond * time.Duration(豪秒數))
	}
}

func (g *Game) 顯示可選位子() []int {
	可選位子 := []int{}

	players := g.Players()
	for i := range players {
		還沒被選擇 := !players[i].IsConnected()
		if 還沒被選擇 {
			可選位子 = append(可選位子, players[i].No())
		}
	}
	return 可選位子
}

func (g *Game) 通知線上人數() {
	players := g.Players()

	線上人數 := []bool{}
	for i := range players {
		線上人數 = append(線上人數, players[i].IsConnected())
	}

	g.旁白(TransferData{
		Action: 更新人數,
		Data:   線上人數,
	}, 0)
}

func (g *Game) 通知更換房主(he Player) {
	g.mx.RLock()
	state := g.gameState
	g.mx.RUnlock()

	g.旁白有話對單個玩家說(he, TransferData{
		Action: 更換房主,
		Data: map[string]bool{
			"遊戲開始": state == 開始階段,
		},
	}, 0)
}

func (g *Game) 加入玩家(he Player) {
	heIsRoomMaster := false
	g.mx.Lock()
	if g.roomMaster == nil {
		g.roomMaster = he
		heIsRoomMaster = true
	}
	g.mx.Unlock()

	if heIsRoomMaster {
		g.通知更換房主(he)
	}

	g.通知線上人數()
}

func (g *Game) 踢除玩家(he Player) {
	if g.是房主(he) {
		var roomMaster Player

		g.mx.Lock()
		g.roomMaster = nil
		for i := range g.players {
			p := g.players[i]
			if p.IsConnected() {
				roomMaster = p
				break
			}
		}
		g.roomMaster = roomMaster
		g.mx.Unlock()

		if roomMaster == nil {
			g.重置()
		} else {
			g.通知更換房主(roomMaster)
		}
	}
	g.通知線上人數()
}

func (g *Game) 是房主(he Player) bool {
	g.mx.RLock()
	yes := he != nil && g.roomMaster != nil && g.roomMaster.No() == he.No()
	g.mx.RUnlock()

	return yes
}
