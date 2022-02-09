package redislock

// 锁对象
type RedSync struct {
	pool RedSyncPool
}

func New(pool RedSyncPool) *RedSync {
	return &RedSync{pool: pool}
}

func (s *RedSync) NewLock(key string, options ...option) *Lock {
	lock := &Lock{
		key:          key,
		identifyFunc: generateIdentity,
		pool:         s.pool,
		stat:         lockStatInit,
		lockTimeout:  defaultLockTimeout,
	}

	for _, opt := range options {
		opt.apply(lock)
	}

	return lock
}
