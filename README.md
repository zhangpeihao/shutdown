# Gracefully shutdown using context.Context and sync.WaitGroup

[![GoDoc](https://godoc.org/github.com/zhangpeihao/shutdown?status.svg)](https://godoc.org/github.com/zhangpeihao/shutdown)
[![Go Report Card](https://goreportcard.com/badge/github.com/zhangpeihao/shutdown)](https://goreportcard.com/report/github.com/zhangpeihao/shutdown)

## Example

```
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
		
		// Some close processes
	}(ctx)

	// Wait interrupt signal and shutdown gracefully
	if err := WaitAndShutdown(ctx, time.Second*5, func(timeout time.Duration) error {
		log.Println("close")
		return nil
	}); err != nil {
		log.Println("Shutdown error:", err)
		return
	}
```