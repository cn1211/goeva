package exp

import "github.com/pkg/errors"

func CheckErr(err error) {
	if err != nil {
		Panic(err)
	}
}

func Panic(err error) {
	if IsPkgErr(err) {
		panic(err)
	}
	panic(errors.WithStack(err))
}

func PanicIf(ok bool, err error) {
	if ok {
		Panic(err)
	}
}

func PanicIfNot(ok bool, err error) {
	if !ok {
		Panic(err)
	}
}

func PanicIfFunc(okFunc func() bool, err error) {
	PanicIf(okFunc(), err)
}
