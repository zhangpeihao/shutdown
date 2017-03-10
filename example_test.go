package shutdown

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"
)

// This example demonstrates the use of shutdown service gracefully
func ExampleShutdown() {
	// Generate a new context
	ctx := NewContext()

	// Run service with this context
	go func(ctx context.Context) {
		if err := ExitWaitGroupAdd(ctx, 1); err != nil {
			return
		}
		defer ExitWaitGroupDone(ctx)

		otherEvent := make(chan struct{})
	FOR_LOOP:
		for {
			select {
			case <-ctx.Done():
				break FOR_LOOP
			case <-otherEvent:
				// ...
			}
		}
	}(ctx)

	// Wait interrupt signal
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	<-sig
	if err := Shutdown(ctx, time.Second*5); err != nil {
		log.Println("Shutdown error:", err)
		return
	}
}
