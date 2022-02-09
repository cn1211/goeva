package redislock

import (
	"time"
)

// 计算下次自旋锁周期
func nextSpinPeriod(retryCdt uint) time.Duration {
	if retryCdt <= 1 {
		return MinObtainSpin
	}
	if retryCdt >= 10 {
		return MaxObtainSpin
	}
	return time.Duration(retryCdt*retryCdt) * time.Millisecond
}

func Args(script *Script, spec string, keysAndArgs []interface{}) []interface{} {
	var args []interface{}
	if script.KeyCount < 0 {
		args = make([]interface{}, 1+len(keysAndArgs))
		args[0] = spec
		copy(args[1:], keysAndArgs)
	} else {
		args = make([]interface{}, 2+len(keysAndArgs))
		args[0] = spec
		args[1] = script.KeyCount
		copy(args[2:], keysAndArgs)
	}
	return args
}
