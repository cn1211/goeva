package retry

import (
	"github.com/cn1211/goeva/core/errorx"
)

const defaultRetryTimes = 3

type (
	RetryOption func(*retryOptions)

	retryOptions struct {
		times int
	}
)

func DoWithRetry(fn func() error, opts ...RetryOption) error {
	options := newRetryOptions()
	for _, opt := range opts {
		opt(options)
	}

	var berr errorx.BatchError
	for i := 0; i < options.times; i++ {
		if err := fn(); err != nil {
			berr.Add(err)
		} else {
			return nil
		}
	}

	return berr.Err()
}

func WithRetry(times int) RetryOption {
	return func(options *retryOptions) {
		options.times = times
	}
}

func newRetryOptions() *retryOptions {
	return &retryOptions{
		times: defaultRetryTimes,
	}
}
