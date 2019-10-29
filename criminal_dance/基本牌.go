package dance

import "sync"

// NewBasicCard 新增第一發現者
func NewBasicCard(owner Player) *BasicCard {
	return &BasicCard{
		owner: owner,
	}
}

// BasicCard 基本牌
type BasicCard struct {
	owner Player
	mx    sync.RWMutex
}

// CanUse 可以使用?
func (c *BasicCard) CanUse() bool {
	owner := c.Owner()
	if owner == nil {
		return false
	}

	return !owner.HasFirstFinder()
}

// ChangeOwner 變更擁有者
func (c *BasicCard) ChangeOwner(player Player) {
	c.mx.Lock()
	c.owner = player
	c.mx.Unlock()
}

// Owner 擁有者
func (c *BasicCard) Owner() (player Player) {
	c.mx.RLock()
	player = c.owner
	c.mx.RUnlock()
	return
}
