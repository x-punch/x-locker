package locker_test

import (
	"sync"
	"testing"
	"time"

	locker "github.com/x-punch/x-locker"
)

var l = locker.NewLock()

func TestLock(t *testing.T) {
	wg, n := sync.WaitGroup{}, 3000
	wg.Add(n)
	for i := 0; i < n; i++ {
		id := "sameid"
		go func() {
			l.Lock(id)
			time.Sleep(time.Millisecond)
			wg.Done()
			l.Unlock(id)
		}()
	}
	wg.Wait()
}

func TestLockEmptyID(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()
	l.Lock("")
}

func TestUnlockEmptyID(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()
	l.Unlock("")
}

func TestUnlockNonexistID(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()
	l.Unlock("id")
}
