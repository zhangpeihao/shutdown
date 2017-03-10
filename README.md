# Gracefully shutdown using context.Context and sync.WaitGroup

[![Build Status](https://travis-ci.org/zhangpeihao/shutdown.svg?branch=master)](https://travis-ci.org/zhangpeihao/shutdown) 
[![Coverage Status](https://coveralls.io/repos/github/zhangpeihao/shutdown/badge.svg?branch=master)](https://coveralls.io/github/zhangpeihao/shutdown?branch=master)
[![GoDoc](https://godoc.org/github.com/zhangpeihao/shutdown?status.svg)](https://godoc.org/github.com/zhangpeihao/shutdown)
[![Go Report Card](https://goreportcard.com/badge/github.com/zhangpeihao/shutdown)](https://goreportcard.com/report/github.com/zhangpeihao/shutdown)
[![Code Climate](https://codeclimate.com/github/zhangpeihao/shutdown/badges/gpa.svg)](https://codeclimate.com/github/zhangpeihao/shutdown)
[![Image Size](https://images.microbadger.com/badges/image/zhangpeihao/shutdown.svg)](https://microbadger.com/images/zhangpeihao/shutdown "Get your own image badge on microbadger.com")

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
	}(ctx)

	// Wait interrupt signal
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	<-sig
	if err := Shutdown(ctx, time.Second*5); err != nil {
		log.Println("Shutdown error:", err)
		return
	}
```