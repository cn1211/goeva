package redislock

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Lock struct {
	key  string   // 锁名称
	stat lockStat // 当前锁状态,用于处理状态同步。

	identifyFunc func() (string, error) // 身份认证生成器
	identify     string                 // 身份值

	lockTimeout  time.Duration // 持有锁时间
	retryTimeout time.Duration // 最大重试时间,为空时默认不重试,非堵塞锁。

	pool RedSyncPool // 连接池对象
}

// 加锁
func (l *Lock) Lock() error {
	identify, err := l.identifyFunc()
	if err != nil {
		return err
	}

	l.identify = identify
	l.stat = lockStatLockFailed

	var (
		retryCdt  uint // 重试次数
		retryTime time.Duration
	)
	for {
		if ok, err := l.acquire(l.identify); err != nil {
			return err
		} else if ok {
			// 抢锁成功
			l.stat = lockStatLocked
			return nil
		}

		retryCdt++
		if retryTime >= l.retryTimeout {
			return errors.New("自旋锁抢占失败")
		}

		// 计算下次自旋锁时间
		retryTime = nextSpinPeriod(retryCdt)
		time.Sleep(retryTime)
	}
}

// 解锁
func (l *Lock) Unlock() error {
	if l.stat != lockStatLocked {
		return errors.Errorf("lock unlock failed because in stat(%d): %s", l.stat, l.key)
	}

	// 释放锁
	ok, err := l.release()
	if err != nil {
		l.stat = lockStatUnlockFailed
		return err
	}
	if !ok {
		l.stat = lockStatUnlockFailed
		return errors.Errorf("lock expired before unlock:%s", l.key)
	}

	// 释放成功
	l.stat = lockStatUnlocked
	return nil
}

// 在原有基础上增加过期时间
func (l *Lock) AddTimeout(addSec int64) error {
	if l.stat != lockStatLocked {
		return errors.Errorf("lock unlock failed because in stat(%d): %s", l.stat, l.key)
	}

	ok, err := l.addTimeout(addSec)
	if err != nil {
		return err
	}

	if !ok {
		l.stat = lockStatUnlockFailed
		return errors.Errorf("lock expired before unlock:%s", l.key)
	}

	// 延时成功
	return nil
}

func (l *Lock) acquire(identity string) (bool, error) {
	conn := l.pool.Get()
	defer conn.Close()

	reply, err := conn.SetNx(l.key, identity, l.lockTimeout)
	if err != nil {
		return false, err
	}

	l.stat = lockStatLocked
	return reply, nil
}

func (l *Lock) release() (bool, error) {
	conn := l.pool.Get()
	defer conn.Close()

	status, err := conn.Eval(releaseScript, l.key, l.identify)
	if err != nil {
		return false, err
	}
	return status != 0, nil
}

// 添加过期时间
func (l *Lock) addTimeout(addSec int64) (bool, error) {
	conn := l.pool.Get()
	defer conn.Close()

	status, err := conn.Eval(addTTLScript, l.key, addSec, l.identify)
	if err != nil {
		return false, err
	}
	return status != 0, nil
}

// 生成身份标识符
// TODO  算法待优化, 后期改用雪花id方案。 https://github.com/sony/sonyflake
func generateIdentity() (string, error) {
	u, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
