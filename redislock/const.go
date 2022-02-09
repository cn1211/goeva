package redislock

import "time"

const (
	// 最大争抢锁持续时长: 30s
	MaxObtainLock = 30 * time.Second
	// 默认持有锁超时时长: 15s
	defaultLockTimeout = 15 * time.Second
	// 最短自旋时长: 4ms
	MinObtainSpin = 4 * time.Millisecond
	// 最大自旋时长: 100ms
	MaxObtainSpin = 100 * time.Millisecond
)

// 锁的状态
type lockStat uint

const (
	// 初始状态
	lockStatInit lockStat = iota
	// 成功得到锁
	lockStatLocked
	// 抢锁失败
	lockStatLockFailed
	// 成功解锁
	lockStatUnlocked
	// 解锁失败
	lockStatUnlockFailed
)
