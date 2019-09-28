package werewolf

import (
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

// NewHuman 建立新玩家
func NewHuman(遊戲 *Game, 位子 int) *Human {
	return &Human{
		位子: 位子,
		遊戲: 遊戲,
	}
}

// Human 玩家
type Human struct {
	位子    int
	開眼睛中  bool
	被投票狀態 bool
	出局狀態  bool
	conn  *websocket.Conn
	遊戲    *Game
	傳話筒   chan 傳輸資料
	讀寫鎖   sync.RWMutex
}

func (我 *Human) 號碼() int {
	return 我.位子
}

func (我 *Human) 閉眼睛() {
	我.開眼睛中 = false
}

func (我 *Human) 開眼睛() {
	我.開眼睛中 = true
}

func (我 *Human) 投票() int {
	if 我.連線() == nil {
		return 0
	}

	var 投票號碼 int
	for {
		so, err := waitChannelBack(我.傳話筒, 選擇玩家)
		if err != nil {
			return 0
		}

		投票號碼, err = strconv.Atoi(so.Reply)
		if err != nil {
			continue
		}

		break
	}
	return 投票號碼
}

func (我 *Human) 被投票(是 bool) {
	我.被投票狀態 = 是
}

func (我 *Human) 被投票了() bool {
	return 我.被投票狀態
}

func (我 *Human) 出局(殺法 KILL) {
	我.出局狀態 = true
	return
}

func (我 *Human) 出局了() bool {
	return 我.出局狀態
}

func (我 *Human) 種族() GROUP {
	return 人質
}

func (我 *Human) 職業() RULE {
	return 平民
}

func (我 *Human) 換位子(新位子 int) int {
	舊的位子 := 我.位子
	我.位子 = 新位子
	return 舊的位子
}

func (我 *Human) 加入(連線 *websocket.Conn) (加入成功 bool) {
	我.讀寫鎖.Lock()
	defer 我.讀寫鎖.Unlock()

	if 我.conn != nil {
		return
	}

	我.conn = 連線
	我.傳話筒 = make(chan 傳輸資料)
	加入成功 = true

	return
}

func (我 *Human) 等待中() {
	我.讀寫鎖.Lock()
	if 我.conn == nil {
		我.讀寫鎖.Unlock()
		return
	}
	我.讀寫鎖.Unlock()

	for {
		so, err := waitSocketBack(我.conn, 無)
		if err != nil {
			我.退出()
			return
		}

		if 我.遊戲.目前階段() == 準備階段 {
			if 我.遊戲.是房主(我) {
				if so.Reply == "start" {
					go func() {
						我.遊戲.開始()
					}()
				}
			}
			continue
		}

		我.傳話筒 <- so
	}
}

func (我 *Human) 退出() {
	我.讀寫鎖.Lock()
	我.conn = nil
	if 我.傳話筒 != nil {
		close(我.傳話筒)
		我.傳話筒 = nil
	}
	我.讀寫鎖.Unlock()
	我.遊戲.踢除玩家(我)
}

func (我 *Human) 已經被選擇() bool {
	我.讀寫鎖.Lock()
	被選擇 := 我.conn != nil
	我.讀寫鎖.Unlock()
	return 被選擇
}

func (我 *Human) 發言() bool {
	我.遊戲.旁白有話對單個玩家說(我, 傳輸資料{
		Display: "請發言",
		Action:  等待回應,
	})
	waitChannelBack(我.傳話筒, 等待回應)
	return false
}

func (我 *Human) 連線() *websocket.Conn {
	我.讀寫鎖.Lock()
	conn := 我.conn
	我.讀寫鎖.Unlock()
	return conn
}

func (我 *Human) 發表遺言() {
	我.遊戲.旁白有話對單個玩家說(我, 傳輸資料{
		Display: "請發言",
		Action:  等待回應,
	})
	waitChannelBack(我.傳話筒, 等待回應)
	return
}
