package dance

import (
	"encoding/json"
	"errors"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

// NewBasicPlayer 新增玩家
func NewBasicPlayer(no int, game *Game) *BasicPlayer {
	name := strconv.Itoa(no) + "號玩家"
	me := &BasicPlayer{
		name:  name,
		no:    no,
		cards: []Card{},
		game:  game,
	}

	return me
}

// BasicPlayer 玩家
type BasicPlayer struct {
	name         string
	no           int
	cards        []Card
	throwCards   []Card
	isCriminal   bool
	isAccomplice bool
	isDetective  bool
	point        int
	conn         *websocket.Conn
	game         *Game
	readCh       chan TransferData
	writeCh      chan TransferData
	mx           sync.RWMutex
}

// Name 號碼
func (me *BasicPlayer) Name() string {
	return me.name
}

// No 號碼
func (me *BasicPlayer) No() int {
	return me.no
}

// Join 加入
func (me *BasicPlayer) Join(conn *websocket.Conn, name string) (ok bool) {
	me.mx.Lock()
	if me.conn == nil {
		ok = true
		me.conn = conn
		me.readCh = make(chan TransferData, 4)
		me.writeCh = make(chan TransferData, 4)
		if name = strings.TrimSpace(name); name != "" {
			me.name = name
		}
	}
	me.mx.Unlock()
	return
}

// Conn 連線
func (me *BasicPlayer) Conn() *websocket.Conn {
	me.mx.Lock()
	conn := me.conn
	me.mx.Unlock()
	return conn
}

// IsConnected 是否連線了
func (me *BasicPlayer) IsConnected() bool {
	me.mx.Lock()
	isConnected := me.conn != nil
	me.mx.Unlock()
	return isConnected
}

// Waiting 等待中
func (me *BasicPlayer) Waiting() {
	if !me.IsConnected() {
		return
	}

	go func() {
		defer func() {
			recover()
		}()
		for {
			data := <-me.writeCh
			err := me.conn.WriteJSON(data)
			if err != nil {
				me.Exit()
				return
			}
		}
	}()
	runtime.Gosched()

	defer func() {
		recover()
	}()

	for {
		so, err := waitSocketBack(me.conn, 無)
		if err != nil {
			me.Exit()
			return
		}

		if me.game.目前階段() == 準備階段 {
			if me.game.是房主(me) {
				td := TransferData{}
				err = json.Unmarshal([]byte(so.Reply), &td)
				if err == nil && td.Action == 更換房主 {
					if td.Reply == "start" {
						go me.game.開始()
						runtime.Gosched()
					}
				}
			}
			continue
		}

		me.readCh <- so
	}
}

// Cards 手上牌組
func (me *BasicPlayer) Cards() []Card {
	tmpCards := []Card{}
	me.mx.Lock()
	for _, card := range me.cards {
		if card.Owner().No() == me.no {
			tmpCards = append(tmpCards, card)
		}
	}
	me.cards = tmpCards
	me.mx.Unlock()
	return tmpCards
}

// ThrowCards 丟出的牌組
func (me *BasicPlayer) ThrowCards() []Card {
	me.mx.Lock()
	cards := me.throwCards
	me.mx.Unlock()
	return cards
}

// IsEmptyCard 牌組是否空了
func (me *BasicPlayer) IsEmptyCard() bool {
	cards := me.Cards()
	me.mx.RLock()
	isEmpty := len(cards) == 0
	me.mx.RUnlock()
	return isEmpty
}

// TurnMe 輪到我
func (me *BasicPlayer) TurnMe(no int) Card {
	for {
		card := me.PlayCard(nil, me.No() == no)
		if card == nil {
			return nil
		}

		if card != nil && card.Name() != 神犬 {
			me.mx.Lock()
			me.throwCards = append(me.throwCards, card)
			me.mx.Unlock()
		}
		return card
	}
}

// PlayCard 出牌
func (me *BasicPlayer) PlayCard(he Player, isTurnMe bool) Card {
	var cards []Card

	if he == nil {
		cards = me.Cards()
	} else {
		cards = he.Cards()
	}

	if len(cards) == 0 {
		return nil
	}

	var targetAction 動作
	if he == nil {
		me.TransferToMe(TransferData{
			Action: 出牌,
			Data:   CardsInfoOutput(cards),
		})
		targetAction = 出牌
	} else {
		cardIndex := []int{}
		for i := range cards {
			cardIndex = append(cardIndex, i)
		}
		me.TransferToMe(TransferData{
			Action: 抽牌,
			Data:   cardIndex,
		})
		targetAction = 抽牌
	}

	var card Card
	var no int
	for {
		so, err := waitChannelBack(me.readCh, targetAction)
		if err != nil {
			no = random(len(cards)) - 1
			if !cards[no].CanUse() && isTurnMe {
				continue
			}
			break
		} else {
			no, err = strconv.Atoi(so.Reply)
			if err == nil && no >= 0 && no < len(cards) {
				break
			}
		}
	}

	card = cards[no]
	if he == nil {
		if card.Name() == 犯人 {
			me.BecomeCriminal(false)
		}

		me.mx.Lock()
		me.cards = append(me.cards[:no], me.cards[no+1:]...)
		Shuffle(me.cards)
		me.mx.Unlock()
	} else if card.Name() == 犯人 {
		he.BecomeCriminal(false)
		// outSideCards = append(outSideCards[:no], outSideCards[no+1:]...)
	}

	return card
}

// TakeCard 拿卡片
func (me *BasicPlayer) TakeCard(card Card) {
	if card == nil {
		return
	}

	if card.Name() == 犯人 {
		me.BecomeCriminal(true)
	}

	card.ChangeOwner(me)
	me.mx.Lock()
	me.cards = append(me.cards, card)
	Shuffle(me.cards)
	me.mx.Unlock()
}

// ClearCard 清空牌
func (me *BasicPlayer) ClearCard() {
	me.mx.Lock()
	me.cards = []Card{}
	me.throwCards = []Card{}
	me.mx.Unlock()
	return
}

// BecomeCriminal 變成犯人
func (me *BasicPlayer) BecomeCriminal(becomeCriminal bool) {
	me.mx.Lock()
	me.isCriminal = becomeCriminal
	me.mx.Unlock()
}

// BecomeAccomplice 變成犯人
func (me *BasicPlayer) BecomeAccomplice(becomeAccomplice bool) {
	me.mx.Lock()
	me.isAccomplice = becomeAccomplice
	me.mx.Unlock()
}

// BecomeDetective 變成偵探
func (me *BasicPlayer) BecomeDetective(becomeDetective bool) {
	me.mx.Lock()
	me.isDetective = becomeDetective
	me.mx.Unlock()
}

// BeAskedCriminal 被問是不是犯人
func (me *BasicPlayer) BeAskedCriminal() bool {

	me.mx.Lock()
	isCriminal := me.isCriminal
	me.mx.Unlock()

	if isCriminal {
		cards := me.Cards()
		for _, card := range cards {
			if card.Name() == 不在場證明 {
				return false
			}
		}

		return true
	}

	return false
}

// IsCriminal 是犯人嗎？
func (me *BasicPlayer) IsCriminal() bool {
	me.mx.RLock()
	isCriminal := me.isCriminal || me.isAccomplice
	me.mx.RUnlock()
	return isCriminal
}

// IsDetective 是偵探嗎？
func (me *BasicPlayer) IsDetective() bool {
	me.mx.RLock()
	isDetective := me.isDetective
	me.mx.RUnlock()
	return isDetective
}

// HasFirstFinder 有第一發現者嗎
func (me *BasicPlayer) HasFirstFinder() bool {
	cards := me.Cards()
	for _, card := range cards {
		if card.Name() == 第一發現者 {
			return true
		}
	}

	return false
}

// TakePoint 拿到分數
func (me *BasicPlayer) TakePoint(p int) {
	me.mx.Lock()
	me.point += p
	me.mx.Unlock()
}

// CurrentPoint 目前分數
func (me *BasicPlayer) CurrentPoint() int {
	me.mx.RLock()
	p := me.point
	me.mx.RUnlock()
	return p
}

// Exit 退出
func (me *BasicPlayer) Exit() {
	me.mx.Lock()
	me.conn = nil
	if me.readCh != nil {
		close(me.readCh)
		me.readCh = nil
	}
	if me.writeCh != nil {
		close(me.writeCh)
		me.writeCh = nil
	}
	me.mx.Unlock()
	me.game.踢除玩家(me)
}

// WaitingAction 等待玩家動作
func (me *BasicPlayer) WaitingAction(targetAction 動作) (TransferData, error) {
	return waitChannelBack(me.readCh, targetAction)
}

// TransferToMe 傳輸資料給我
func (me *BasicPlayer) TransferToMe(data TransferData) error {
	if !me.IsConnected() {
		return errors.New("尚未連線")
	}
	me.writeCh <- data
	runtime.Gosched()
	return nil
}
