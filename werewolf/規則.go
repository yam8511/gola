package werewolf

import (
	"github.com/gorilla/websocket"
)

type 淘汰者 struct {
	殺法 KILL
	玩家 Player
}

type 傳輸資料 struct {
	UID     string      `json:"uid"`
	Sound   string      `json:"sound"`
	Display string      `json:"display"`
	Action  動作          `json:"action"`
	Data    interface{} `json:"data"`
	Reply   string      `json:"reply"`
}

type 動作 string

const (
	無       = 動作("")
	給傳話筒    = 動作("for_channel")
	角色設定    = 動作("role_setup")
	選擇號碼    = 動作("select_number")
	拿到角色    = 動作("take_rule")
	拿到Token = 動作("take_token")
	更換房主    = 動作("change_room_master")
	遊戲已開始   = 動作("game_is_running")
	遊戲結束    = 動作("game_over")
	天黑請閉眼   = 動作("all_close_eyes")
	天亮請睜眼   = 動作("all_open_eyes")
	選擇玩家    = 動作("select_player")
	等待回應    = 動作("waiting")
	顯示投票結果  = 動作("vote_result")
	玩家淘汰    = 動作("player_out")
	更新人數    = 動作("refresh_online")
)

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

type 時間 string

const (
	白天 = 時間("白天")
	投票 = 時間("投票")
	黑夜 = 時間("黑夜")
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
	// http://www.87g.com/lrs/59119.html
	平民  = RULE("平民")
	預言家 = RULE("預言家")
	熊   = RULE("熊")
	女巫  = RULE("女巫")
	獵人  = RULE("獵人")
	騎士  = RULE("騎士")
	狼人  = RULE("狼人")
	狼王  = RULE("狼王")
	雪狼  = RULE("雪狼")
	狼美人 = RULE("狼美人")
)

// KILL 殺法
type KILL string

const (
	票殺 = KILL("票殺")
	狼殺 = KILL("狼殺")
	毒殺 = KILL("毒殺")
	獵殺 = KILL("獵殺")
	騎殺 = KILL("騎殺")
	自爆 = KILL("自爆")
)

// Player 玩家
type Player interface {
	加入(*websocket.Conn) (加入成功 bool)
	等待中()
	退出()
	號碼() int
	閉眼睛()
	開眼睛()
	投票(string) int
	出局(KILL)
	出局了() bool
	種族() GROUP
	職業() RULE
	換位子(int) int
	已經被選擇() bool
	發言(投票發言 bool) (進入黑夜 bool)
	連線() *websocket.Conn
	發表遺言()
	等待動作(動作, string) (傳輸資料, error)
	傳話給玩家(傳輸資料) error
}

// Skiller 有能力的人
type Skiller interface {
	Player
	能力(第幾個技能 int) (用來表示狼人咬誰 Player)
	需要夜晚行動() bool
}
