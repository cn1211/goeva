package redislock

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"time"
)

type RedSyncPool interface {
	Get() Conn
}

// 连接对象方法
type Conn interface {
	SetNx(key string, value interface{}, expireAt time.Duration) (bool, error)
	Eval(script *Script, keysAndArgs ...interface{}) (interface{}, error) // 执行指定脚本
	Close() error                                                         // 关闭redis对象
}

type Script struct {
	KeyCount int
	Src      string
	Hash     string
}

func NewScript(keyCount int, src string) *Script {
	h := sha1.New()
	_, _ = io.WriteString(h, src)
	return &Script{
		KeyCount: keyCount,
		Src:      src,
		Hash:     hex.EncodeToString(h.Sum(nil)),
	}
}
