package dance

import "github.com/gorilla/websocket"

// https://bit.ly/2PbEsBW
// https://punchboardgame.pixnet.net/blog/post/439177456

type 規則設定 struct {
	CombineName string           `json:"CombineName"`
	PlayerCount int              `json:"PlayerCount"`
	TargetPoint int              `json:"TargetPoint"`
	Advanced    bool             `json:"Advanced"`
	CardSet     map[CardName]int `json:"CardSet"`
}

var 規則牌組設定 = map[string]規則設定{
	"3": {
		"3人場",
		3, 5, false, map[CardName]int{
			第一發現者: 1,
			犯人:    1,
			偵探:    1,
			不在場證明: 2,
			目擊者:   1,
			情報交換:  2,
			交易:    1,
			謠言:    2,
			神犬:    1,
		},
	},
	"4": {
		"4人場",
		4, 5, false, map[CardName]int{
			第一發現者: 1,
			犯人:    1,
			偵探:    1,
			不在場證明: 1,
			共犯:    1,
			隨機:    11,
		},
	},
	"5": {
		"5人場",
		5, 10, false, map[CardName]int{
			第一發現者: 1,
			犯人:    1,
			偵探:    1,
			不在場證明: 2,
			共犯:    1,
			隨機:    14,
		},
	},
	"6": {
		"6人場",
		6, 10, false, map[CardName]int{
			第一發現者: 1,
			犯人:    1,
			偵探:    2,
			不在場證明: 2,
			共犯:    2,
			隨機:    16,
		},
	},
	"7": {
		"7人場",
		7, 10, false, map[CardName]int{
			第一發現者: 1,
			犯人:    1,
			偵探:    2,
			不在場證明: 3,
			共犯:    2,
			隨機:    19,
		},
	},
	"8-basic": {
		"8人基本場",
		8, 10, false, map[CardName]int{
			第一發現者: 1,
			偵探:    4,
			犯人:    1,
			共犯:    2,
			普通人:   2,
			目擊者:   3,
			不在場證明: 5,
			情報交換:  4,
			交易:    4,
			謠言:    5,
			神犬:    1,
		},
	},
	// "8-advanced": {
	// 	"8人進階場",
	// 	8, 10, true, map[CardName]int{
	// 		第一發現者: 1,
	// 		偵探:    4,
	// 		犯人:    1,
	// 		共犯:    2,
	// 		普通人:   2,
	// 		目擊者:   2,
	// 		不在場證明: 5,
	// 		情報交換:  4,
	// 		交易:    4,
	// 		謠言:    4,
	// 		神犬:    1,
	// 		警部:    1,
	// 		少年:    1,
	// 	},
	// },
	// "8-random": {
	// 	"8人隨機場",
	// 	8, 10, true, map[CardName]int{
	// 		第一發現者: 1,
	// 		偵探:    4,
	// 		犯人:    1,
	// 		共犯:    2,
	// 		普通人:   2,
	// 		目擊者:   3,
	// 		隨機:    19,
	// 	},
	// },
}

var 基本牌組 = map[CardName]int{
	第一發現者: 1,
	偵探:    4,
	犯人:    1,
	共犯:    2,
	普通人:   2,
	目擊者:   3,
	不在場證明: 5,
	情報交換:  4,
	交易:    4,
	謠言:    5,
	神犬:    1,
	// 警部:    1,
	// 少年:    1,
}

// CardName 卡片名稱
type CardName string

const (
	隨機    = CardName("隨機")
	第一發現者 = CardName("第一發現者")
	偵探    = CardName("偵探")
	犯人    = CardName("犯人")
	共犯    = CardName("共犯")
	普通人   = CardName("普通人")
	目擊者   = CardName("目擊者")
	不在場證明 = CardName("不在場證明")
	交易    = CardName("交易")
	情報交換  = CardName("情報交換")
	謠言    = CardName("謠言")
	神犬    = CardName("神犬")
	警部    = CardName("警部")
	少年    = CardName("少年")
)

// GameResult 遊戲結果
type GameResult string

const (
	進行中  = GameResult("進行中")
	偵探勝利 = GameResult("偵探勝利")
	神犬勝利 = GameResult("神犬勝利")
	犯人勝利 = GameResult("犯人勝利")
)

type 遊戲階段 int

const (
	初始階段 = 遊戲階段(iota)
	準備階段
	開始階段
	結算階段
)

type 動作 string

const (
	無       = 動作("")
	給傳話筒    = 動作("for_channel")
	牌組設定    = 動作("card_setup")
	選擇號碼    = 動作("select_number")
	拿到角色    = 動作("take_rule")
	拿到Token = 動作("take_token")
	更換房主    = 動作("change_room_master")
	遊戲已開始   = 動作("game_is_running")
	遊戲結束    = 動作("game_over")
	輪到我     = 動作("turn_me")
	出牌      = 動作("play_card")
	顯示出牌    = 動作("show_play_card")
	抽牌      = 動作("draw_card")
	顯示被抽牌   = 動作("show_draw_card")
	顯示拿牌    = 動作("show_take_card")
	顯示結算    = 動作("show_result")
	等待回應    = 動作("waiting")
	更新人數    = 動作("refresh_online")
	更新桌面    = 動作("refresh_desktop")
	看手牌     = 動作("look_player_cards")
)

// TransferData 傳輸資料
type TransferData struct {
	Sound   string      `json:"sound"`
	Display string      `json:"display"`
	Action  動作          `json:"action"`
	Data    interface{} `json:"data"`
	Reply   string      `json:"reply"`
}

// CardDisplay 卡片顯示
type CardDisplay struct {
	Index   int      `json:"Index"`
	Name    CardName `json:"Name"`
	Detail  string   `json:"Detail"`
	Disable bool     `json:"Disable"`
}

// Card 卡片
type Card interface {
	Name() CardName
	Detail() string
	Skill(*Game) GameResult
	CanUse() bool
	ChangeOwner(Player)
	Owner() Player
}

// Player 玩家
type Player interface {
	Name() string
	No() int
	Join(*websocket.Conn, string) (加入成功 bool)
	Conn() *websocket.Conn
	TakePoint(int)
	CurrentPoint() int
	IsConnected() bool
	Waiting()
	Cards() []Card
	ThrowCards() []Card
	IsEmptyCard() bool
	HasFirstFinder() bool
	TurnMe(int) Card
	PlayCard(he Player, isTurnMe bool) Card
	TakeCard(Card)
	ClearCard()
	BecomeCriminal(becomCrimnal bool)
	BecomeAccomplice(becomeAccomplice bool)
	BecomeDetective(becomeDetective bool)
	IsCriminal() bool
	IsDetective() bool
	BeAskedCriminal() bool
	Exit()
	WaitingAction(動作) (TransferData, error)
	TransferToMe(TransferData) error
}
