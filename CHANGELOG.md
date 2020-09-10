### 1.2.3 (2020/09/10)
> 1. Refact lock with interface.
> 2. Refact errors.

### 1.2.2 (2020/09/04)
> 1. Fixing nil error issue: If SET NX command failed with redis.Nil, means key already locked, should return acquire lock failed error.

### 1.2.1 (2020/08/27)
> 1. fixing interface missing NewLock method.

### 1.2.0 (2020/08/27)
> 1. Change to expose interface.

### 1.1.0 (2020/07/03)
> 1. add redlock for Distributed Lock.

### 1.0.0 (2020/05/25)
> 1. finish group locker.