package redlock

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var _ Redlock = &redlock{}

// A DelayFunc is used to decide the amount of time to wait between retries.
type DelayFunc func(tries int) time.Duration

// A Redlock is a distributed mutual exclusion lock.
type Redlock interface {
	Lock() error
	Unlock() (bool, error)
	Extend() (bool, error)
	Valid() (bool, error)
}

type redlock struct {
	id         string
	expiration time.Duration

	maxTries  int
	delayFunc DelayFunc

	factor float64

	quorum int

	genValueFunc func() (string, error)
	value        string
	until        time.Time

	clients []RedisClient
}

// Lock will lock the id. In case it returns an error on failure, you may retry to acquire the lock by calling this method again.
func (l *redlock) Lock() error {
	value, err := l.genValueFunc()
	if err != nil {
		return err
	}

	for i := 0; i < l.maxTries; i++ {
		if i != 0 {
			time.Sleep(l.delayFunc(i))
		}

		start := time.Now()
		n, err := l.actOnClientsAsync(func(client RedisClient) (bool, error) {
			return l.acquire(client, value)
		})
		if n == 0 && err != nil {
			return err
		}

		now := time.Now()
		until := now.Add(l.expiration - now.Sub(start) - time.Duration(int64(float64(l.expiration)*l.factor)))
		if n >= l.quorum && now.Before(until) {
			l.value = value
			l.until = until
			return nil
		}
		l.actOnClientsAsync(func(client RedisClient) (bool, error) {
			return l.release(client, value)
		})
	}

	return ErrAcquireLockFailed
}

// Unlock unlocks m and returns the status of unlock.
func (l *redlock) Unlock() (bool, error) {
	n, err := l.actOnClientsAsync(func(client RedisClient) (bool, error) {
		return l.release(client, l.value)
	})
	if n < l.quorum {
		return false, err
	}
	return true, nil
}

// Extend resets the mutex's expiry and returns the status of expiry extension.
func (l *redlock) Extend() (bool, error) {
	n, err := l.actOnClientsAsync(func(client RedisClient) (bool, error) {
		return l.extend(client, l.value, int(l.expiration/time.Millisecond))
	})
	if n < l.quorum {
		return false, err
	}
	return true, nil
}

// Valid will check liveness for lock
func (l *redlock) Valid() (bool, error) {
	n, err := l.actOnClientsAsync(func(client RedisClient) (bool, error) {
		return l.valid(client)
	})
	return n >= l.quorum, err
}

func (l *redlock) valid(client RedisClient) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), l.expiration)
	defer cancel()
	result, err := client.Do(ctx, "GET", l.id).Text()
	if err != nil {
		return false, err
	}
	return result == l.value, nil
}

func (l *redlock) acquire(client RedisClient, value string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), l.expiration)
	defer cancel()
	result, err := client.Do(ctx, "SET", l.id, value, "NX", "PX", int(l.expiration/time.Millisecond)).Text()
	if err != nil {
		if err == redis.Nil {
			return false, ErrAcquireLockFailed
		}
		return false, err
	}
	return result == "OK", nil
}

var deleteScript = redis.NewScript(`
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end
`)

func (l *redlock) release(client RedisClient, value string) (bool, error) {
	status, err := deleteScript.Run(context.TODO(), client, []string{l.id}, value).Int64()
	return err == nil && status != 0, err
}

var expireScript = redis.NewScript(`
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("PEXPIRE", KEYS[1], ARGV[2])
	else
		return 0
	end
`)

func (l *redlock) extend(client RedisClient, value string, expiry int) (bool, error) {
	status, err := expireScript.Run(context.TODO(), client, []string{l.id}, value, expiry).Int64()
	return err == nil && status != 0, err
}

func (l *redlock) actOnClientsAsync(actFn func(RedisClient) (bool, error)) (int, error) {
	type result struct {
		Success bool
		Err     error
	}

	ch := make(chan result)
	for _, client := range l.clients {
		go func(client RedisClient) {
			r := result{}
			r.Success, r.Err = actFn(client)
			ch <- r
		}(client)
	}
	n := 0
	err := &Error{}
	for range l.clients {
		r := <-ch
		if r.Success {
			n++
		} else if r.Err != nil {
			err.Append(r.Err)
		}
	}
	return n, err
}
