# X Locker
sync.Mutex can be used to lock in golang, but you need to define the mutex in advance.
Sometimes what we want to lock is generated dynamicly, so this package is used to lock in dynamic.
We can group some lock by id, and they shared the same mutex.

# Usage
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
