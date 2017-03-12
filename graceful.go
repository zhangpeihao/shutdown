package shutdown

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

// ContextKey context key type
type ContextKey string

const (
	// ContextCancel cancel func value name in context
	ContextCancel ContextKey = "shutdown/cancel"
	// ContextExitWaitGroup waitgroup object name in context
	ContextExitWaitGroup ContextKey = "shutdown/exitWaitGroup"
)

var (
	// ErrNoCancel no cancel value error
	ErrNoCancel = errors.New("no cancel value in context")
	// ErrCancelValueTypeError cancel value in context type error
	ErrCancelValueTypeError = errors.New("cancel value in context type error")
	// ErrNoExitWaitGroup no exitWaitGroup value in context
	ErrNoExitWaitGroup = errors.New("no exitWaitGroup value in context")
	// ErrExitWaitGroupValueTypeError exitWaitGroup value in context type error
	ErrExitWaitGroupValueTypeError = errors.New("exitWaitGroup value in context type error")
	// ErrShutdownTimeout shutdown timeout
	ErrShutdownTimeout = errors.New("shutdown timeout")
)

// NewContext create new context
func NewContext() (ctx context.Context) {
	var exitWaitGroup sync.WaitGroup
	var ctxBase context.Context
	var cancel context.CancelFunc
	ctxBase, cancel = context.WithCancel(context.Background())
	ctx = context.WithValue(context.WithValue(ctxBase, ContextExitWaitGroup, &exitWaitGroup), ContextCancel, cancel)
	return
}

// Shutdown shutdown gracefully
func Shutdown(ctx context.Context, shutdownTimeout time.Duration) error {
	value := ctx.Value(ContextCancel)
	if value == nil {
		return ErrNoCancel
	}
	cancel, ok := value.(context.CancelFunc)
	if !ok {
		return ErrCancelValueTypeError
	}
	value = ctx.Value(ContextExitWaitGroup)
	if value == nil {
		return ErrNoExitWaitGroup
	}
	exitWaitGroup, ok := value.(*sync.WaitGroup)
	if !ok {
		return ErrExitWaitGroupValueTypeError
	}
	exitSignal := make(chan struct{})
	go func() {
		cancel()
		exitWaitGroup.Wait()
		exitSignal <- struct{}{}
	}()
	select {
	case <-exitSignal:
		fmt.Println("Exit successfully")
	case <-time.After(shutdownTimeout):
		return ErrShutdownTimeout
	}
	if err := ctx.Err(); err != context.Canceled {
		return err
	}
	return nil
}

// WaitAndShutdown shutdown gracefully
func WaitAndShutdown(ctx context.Context, shutdownTimeout time.Duration) error {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	<-sig
	return Shutdown(ctx, shutdownTimeout)
}

// ExitWaitGroupAdd waitgroup counter adds delta
func ExitWaitGroupAdd(ctx context.Context, i int) error {
	value := ctx.Value(ContextExitWaitGroup)
	if value == nil {
		return ErrNoExitWaitGroup
	}
	exitWaitGroup, ok := value.(*sync.WaitGroup)
	if !ok {
		return ErrExitWaitGroupValueTypeError
	}
	exitWaitGroup.Add(i)
	return nil
}

// ExitWaitGroupDone waitgroup down
func ExitWaitGroupDone(ctx context.Context) error {
	value := ctx.Value(ContextExitWaitGroup)
	if value == nil {
		return ErrNoExitWaitGroup
	}
	exitWaitGroup, ok := value.(*sync.WaitGroup)
	if !ok {
		return ErrExitWaitGroupValueTypeError
	}
	exitWaitGroup.Done()
	return nil
}
