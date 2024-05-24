package sigterm

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// CancelContext returns a context that wraps the given context with a cancel
// function that's called when a SIGTERM or SIGINT is trapped.
func CancelContext(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	ctxWithCancel, cancel := context.WithCancel(ctx)

	go func() {
		defer cancel()

		term := make(chan os.Signal, 1)
		signal.Notify(term, syscall.SIGTERM, syscall.SIGINT)

		select {
		case <-term:
		case <-ctx.Done():
		}
	}()

	return ctxWithCancel
}
