package chat

import (
	"sync"
)

type Binder struct {
	userId2ConnMap map[uint64]*Conn
	mu             sync.RWMutex
}

func (b *Binder) Bind(uid uint64, conn *Conn) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.userId2ConnMap[uid] = conn
}

func (b *Binder) UnBind(uid uint64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.userId2ConnMap[uid]; ok {
		delete(b.userId2ConnMap, uid)
	}
}

func (b *Binder) ConnByUid(uid uint64) (*Conn, error) {
	conn, ok := b.userId2ConnMap[uid]
	if !ok {
		return nil, ERR03
	}
	return conn, nil
}
