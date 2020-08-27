package redlock

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

var (
	// DefaultLockExpiration represents lock expire duration
	DefaultLockExpiration = 10 * time.Second
	// DefaultMaxRetries represents default max retries for lock request
	DefaultMaxRetries = 32
	// DefaultDriftFactor represents drift factor
	DefaultDriftFactor = 0.01
)

// Locker provides a simple method for creating distributed locks using multiple Redis servers.
type Locker interface {
}

type locker struct {
	clients []RedisClient
}

// New will creates a Redlock locker
func New(clients []RedisClient) Locker {
	return &locker{
		clients: clients,
	}
}

// NewLock returns a new distributed lock with given id.
func (r *locker) NewLock(id string, options ...Option) Redlock {
	m := &redlock{
		id:           id,
		expiration:   DefaultLockExpiration,
		maxTries:     DefaultMaxRetries,
		factor:       DefaultDriftFactor,
		delayFunc:    func(tries int) time.Duration { return 500 * time.Millisecond },
		genValueFunc: genRandomValue,
		quorum:       len(r.clients)/2 + 1,
		clients:      r.clients,
	}
	for _, o := range options {
		o.Apply(m)
	}
	return m
}

// An Option configures a mutex.
type Option interface {
	Apply(*redlock)
}

// OptionFunc is a function that configures a mutex.
type OptionFunc func(*redlock)

// Apply calls f(mutex)
func (f OptionFunc) Apply(mutex *redlock) {
	f(mutex)
}

// SetExpiry can be used to set the expiry of a mutex to the given value.
func SetExpiry(expiry time.Duration) Option {
	return OptionFunc(func(m *redlock) {
		m.expiration = expiry
	})
}

// SetMaxTries can be used to set the max number of times lock acquire is attempted.
func SetMaxTries(maxTries int) Option {
	return OptionFunc(func(m *redlock) {
		m.maxTries = maxTries
	})
}

// SetRetryDelay can be used to set the amount of time to wait between retries.
func SetRetryDelay(delay time.Duration) Option {
	return OptionFunc(func(m *redlock) {
		m.delayFunc = func(tries int) time.Duration {
			return delay
		}
	})
}

// SetRetryDelayFunc can be used to override default delay behavior.
func SetRetryDelayFunc(delayFunc DelayFunc) Option {
	return OptionFunc(func(m *redlock) {
		m.delayFunc = delayFunc
	})
}

// SetDriftFactor can be used to set the clock drift factor.
func SetDriftFactor(factor float64) Option {
	return OptionFunc(func(m *redlock) {
		m.factor = factor
	})
}

// SetGenValueFunc can be used to set the custom value generator.
func SetGenValueFunc(genValueFunc func() (string, error)) Option {
	return OptionFunc(func(m *redlock) {
		m.genValueFunc = genValueFunc
	})
}

func genRandomValue() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}
