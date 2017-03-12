package shutdown

import (
	"context"
	"log"
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

	// Wait interrupt signal and shutdown gracefully
	if err := WaitAndShutdown(ctx, time.Second*5, func(timeout time.Duration) error {
		log.Println("close")
		return nil
	}); err != nil {
		log.Println("Shutdown error:", err)
		return
	}
}
