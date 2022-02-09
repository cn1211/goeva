package exp

import (
	goerr "errors"

	"github.com/pkg/errors"
)

// src: github.com/pkg/errors/stack.go::*stack
type stackTracer interface {
	StackTrace() errors.StackTrace
}

// src: github.com/pkg/errors/errors.go::Cause(err error) error
type causer interface {
	Cause() error
}

// 判断某个error是否由pkg/errors生成的(错误链中至少包含一个即可)
func IsPkgErr(err error) bool {
	isPkgErr := func(err error) bool {
		_, ok := err.(stackTracer)
		return ok
	}

	// 递归寻找
	for err != nil {
		if isPkgErr(err) {
			return true
		}

		// pkg/errors
		cause, ok := err.(causer)
		if ok {
			err = cause.Cause()
			continue
		}

		// errors
		unpackErr := goerr.Unwrap(err)
		if unpackErr == nil {
			break
		}
		err = unpackErr
	}
	return false
}
