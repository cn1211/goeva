package exp

import (
	"errors"
	"fmt"
	"strings"

	"github.com/facebookgo/stack"
	"github.com/sirupsen/logrus"
)

// Try 协程异常捕获
func Try(fun func()) {
	defer func() {
		if err := recover(); err != nil {
			c := stack.Callers(2)
			for i := len(c) - 1; i >= 0; i-- {
				frame := c[i]
				if frame.Name == "Try" && frame.Line == 29 && strings.HasSuffix(frame.File, "recover.go") {
					c = c[:i+1]
					break
				}
			}

			logrus.Errorf(fmt.Sprintf("panic:%v stacks:%s", err, c.String()))
		}
	}()

	fun()
}

// TryWithErr  协程异常捕获
func TryWithErr(fun func(), handle func(err error)) {
	defer func() {
		if err := recover(); err != nil {
			c := stack.Callers(2)
			for i := len(c) - 1; i >= 0; i-- {
				frame := c[i]
				if frame.Name == "Try" && frame.Line == 51 && strings.HasSuffix(frame.File, "recover.go") {
					c = c[:i+1]
					break
				}
			}

			logrus.Errorf(fmt.Sprintf("panic:%v stacks:%s", err, c.String()))

			handle(errors.New(fmt.Sprint(err)))
		}
	}()

	fun()
}
