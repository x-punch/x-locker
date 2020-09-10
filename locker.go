package locker

import (
	"sync"
)

// Lock represents lock which can lock by string, used in Monolithic application
type Lock interface {
	Lock(string) error
	Unlock(string) error
}

type lock struct {
	mu     sync.Mutex
	groups map[string]*mutexGroup
}

type mutexGroup struct {
	count int
	mu    sync.Mutex
}

// NewLock will create lock instance
func NewLock() Lock {
	return &lock{groups: make(map[string]*mutexGroup)}
}

// Lock will create or get mutex group and then lock it by id
func (l *lock) Lock(id string) error {
	if len(id) == 0 {
		return ErrInvalidLockGroup
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
	return nil
}

// Unlock will get mutex group and unlock it by id
func (l *lock) Unlock(id string) error {
	if len(id) == 0 {
		return ErrInvalidLockGroup
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
		return ErrLockGroupNotFound
	}
	return nil
}
