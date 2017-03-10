package shutdown

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestGraceful(t *testing.T) {
	ctx := NewContext()
	for i := 0; i < 20; i++ {
		go func(ctx context.Context, i int) {
			if err := ExitWaitGroupAdd(ctx, 1); err != nil {
				t.Fatalf("BlockService goroutine %d ExitWaitGroupAdd error: %s", i, err)
				return
			}
			defer ExitWaitGroupDone(ctx)

			fmt.Printf("BlockService goroutine %d running\n", i)
			<-ctx.Done()
			defer fmt.Printf("BlockService goroutine %d done\n", i)
			fmt.Printf("BlockService goroutine %d closing\n", i)
			time.Sleep(time.Second)
		}(ctx, i)
	}
	// Shutdown after 2 seconds
	time.Sleep(time.Second * 2)
	if err := Shutdown(ctx, time.Second*5); err != nil {
		t.Error("GracefulShutdown error:", err)
		return
	}
	log.Println("Shutdown OK")
}
