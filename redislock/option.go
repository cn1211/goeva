package redislock

import "time"

type option interface {
	apply(*Lock)
}

type OptionFunc func(lock *Lock)

func (f OptionFunc) apply(l *Lock) {
	f(l)
}

// 设置自定义身份认证器
func WithIdentityFunc(identifyFunc func() (string, error)) option {
	return OptionFunc(func(l *Lock) {
		l.identifyFunc = identifyFunc
	})
}

// 设置持有锁的超时时间
func WithLockTimeout(lockTimeout time.Duration) option {
	return OptionFunc(func(l *Lock) {
		l.lockTimeout = lockTimeout
	})
}

// 抢锁重试时长
func WithRetryTimeout(retryTimeout time.Duration) option {
	return OptionFunc(func(l *Lock) {
		l.retryTimeout = retryTimeout
	})
}
