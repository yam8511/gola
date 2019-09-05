package main

type 種族 string

const (
	人質 = 種族("人質")
	神職 = 種族("神職")
	狼職 = 種族("狼職")
)

type 玩家 interface {
	號碼() int
	閉眼睛()
	開眼睛()
	投票() (int, error)
	被投票(bool)
	被投票了() bool
	出局了() bool
	種族() 種族
	換號碼(int) int
}

type 能力者 interface {
	玩家
	能力()
}
