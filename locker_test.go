package locker_test

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	locker "github.com/x-punch/x-locker"
)

var l = locker.NewLock()

func TestLockWithSameGroup(t *testing.T) {
	wg, n := sync.WaitGroup{}, 3000
	group := "sameid"
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(g string) {
			l.Lock(g)
			defer l.Unlock(g)
			time.Sleep(time.Millisecond)
			wg.Done()
		}(group)
	}
	wg.Wait()
}

func TestLockWithRandomGroup(t *testing.T) {
	wg, n := sync.WaitGroup{}, 3000
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(g string) {
			l.Lock(g)
			defer l.Unlock(g)
			time.Sleep(time.Millisecond)
			wg.Done()
		}(fmt.Sprintf("%x", rand.Intn(n)))
	}
	wg.Wait()
}

func TestLockEmptyID(t *testing.T) {
	err := l.Lock("")
	if err == nil || !errors.Is(err, locker.ErrInvalidLockGroup) {
		t.Fail()
	}
}

func TestUnlockEmptyID(t *testing.T) {
	err := l.Unlock("")
	if err == nil || !errors.Is(err, locker.ErrInvalidLockGroup) {
		t.Fail()
	}
}

func TestUnlockNonexistID(t *testing.T) {
	err := l.Unlock("id")
	if err == nil || !errors.Is(err, locker.ErrLockGroupNotFound) {
		t.Fail()
	}
}
