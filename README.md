# X Locker
sync.Mutex can be used to lock in golang, but you need to define the mutex in advance.
Sometimes what we want to lock is generated dynamicly, so this package is used to lock in dynamic.
We can group some lock by id, and they shared the same mutex.

## Usage
```go
import locker "github.com/x-punch/x-locker"
```
```go
l := locker.NewLocker()
```
```go
l.Lock("id")
defer l.Unlock("id")
// do something
```

# Redlock
Redsync provides a Redis-based distributed mutual exclusion lock implementation for Go as described in [this post](http://redis.io/topics/distlock).

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

	l := locker.NewLock("id")
	if err := l.Lock(); err != nil {
		panic(err)
	}
	defer l.Unlock()
    
    // do something
}
```