# Redlock

Redsync provides a Redis-based distributed mutual exclusion lock implementation for Go as described in [this post](http://redis.io/topics/distlock).

This lib is inspired by [go-redsync/redsync](github.com/go-redsync/redsync), but redis dependency changed to [go-redis](github.com/go-redis/redis).

## Usage
```go
package main

import (
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/x-punch/x-locker/redlock"
)

func main() {
	locker := redlock.New([]redlock.RedisClient{redis.NewClient(&redis.Options{Addr: ":6379"})})

	l := locker.NewLock("test")
	if err := l.Lock(); err != nil {
		panic(err)
	}
	defer l.Unlock()
    
    // do something
}
```
