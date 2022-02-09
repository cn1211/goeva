package signal

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
)

func WaitForStop(ctx context.Context) error {
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt, os.Kill, syscall.SIGKILL, syscall.SIGTERM)
	select {
	case sig := <-sigCh:
		return errors.New(fmt.Sprintf("接收到关闭信号(%s)", sig.String()))
	case <-ctx.Done():
		return ctx.Err()
	}
}
