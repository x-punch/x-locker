package locker

import (
	"sync"
)

// Lock represents lock which can lock by string, used in Monolithic application
type Lock struct {
	mu     sync.Mutex
	groups map[string]*mutexGroup
}

type mutexGroup struct {
	count int
	mu    sync.Mutex
}

// NewLock will create lock instance
func NewLock() *Lock {
	return &Lock{groups: make(map[string]*mutexGroup)}
}

// Lock will create or get mutex group and then lock it by id
func (l *Lock) Lock(id string) {
	if len(id) == 0 {
		panic("lock id cannot by empty")
	}
	l.mu.Lock()
	g, ok := l.groups[id]
	if !ok {
		g = &mutexGroup{count: 1}
		l.groups[id] = g
	} else {
		g.count++
	}
	l.mu.Unlock()
	g.mu.Lock()
}

// Unlock will get mutex group and unlock it by id
func (l *Lock) Unlock(id string) {
	if len(id) == 0 {
		panic("lock id cannot by empty")
	}
	l.mu.Lock()
	if g, ok := l.groups[id]; ok {
		g.count--
		if g.count == 0 {
			delete(l.groups, id)
		}
		l.mu.Unlock()
		g.mu.Unlock()
	} else {
		panic("unlock id not found")
	}
}
