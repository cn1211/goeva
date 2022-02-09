package httplib

import (
	"net"
	"net/http"
	"time"
)

var (
	defaultClient = NewClient(time.Second * 10)
)

func NewTransport() *http.Transport {
	return &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

func NewClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Transport: NewTransport(),
		Timeout:   timeout,
	}
}
