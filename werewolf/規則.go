package werewolf

import "github.com/gorilla/websocket"

type 階段 int

const (
	準備階段 = 階段(iota)
	開始階段
)

type 遊戲結果 string

const (
	進行中 = 遊戲結果("進行中")
	平手  = 遊戲結果("平手")
	人勝  = 遊戲結果("人勝")
	狼勝  = 遊戲結果("狼勝")
)

// GROUP 種族
type GROUP string

const (
	人質 = GROUP("人質")
	神職 = GROUP("神職")
	狼職 = GROUP("狼職")
)

// RULE 角色
type RULE string

const (
	平民  = RULE("平民")
	預言家 = RULE("預言家")
	女巫  = RULE("女巫")
	獵人  = RULE("獵人")
	騎士  = RULE("騎士")
	狼人  = RULE("狼人")
	狼王  = RULE("狼王")
	雪狼  = RULE("雪狼")
)

// KILL 殺法
type KILL string

const (
	票殺 = KILL("票殺")
	狼殺 = KILL("狼殺")
	毒殺 = KILL("毒殺")
	獵殺 = KILL("獵殺")
	騎殺 = KILL("騎殺")
)

// Player 玩家
type Player interface {
	加入(*websocket.Conn)
	退出()
	號碼() int
	閉眼睛()
	開眼睛()
	投票() int
	出局(KILL)
	出局了() bool
	種族() GROUP
	職業() RULE
	換位子(int) int
	已經被選擇() bool
	發言()
}

// Skiller 有能力的人
type Skiller interface {
	Player
	能力()
	需要夜晚行動() bool
}
