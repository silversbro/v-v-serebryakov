package main

import (
	"context"
	"fmt"
	"time"
)

type RequestID string

func main() {
	fmt.Println("app starting")

	ctx := context.Background()

	cancelCtx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	cancelCtx, cancel := context.WithTimeout(cancelCtx, 10*time.Second)
	defer cancel()

	context.AfterFunc(cancelCtx, func() {
		fmt.Println("context done afterFunc")
	})

	ctx = WithRequestID(cancelCtx, 111)
	ctx = context.WithValue(ctx, "user_id", 222)

	go task(ctx)

	// hard work...
	time.Sleep(5 * time.Second)

	fmt.Println("app end")
}

func WithRequestID(ctx context.Context, val int) context.Context {
	return context.WithValue(ctx, RequestID("request_id"), val)
}

func task(ctx context.Context) {
	//ctx = context.WithoutCancel(ctx)

	ctx, cancel := context.WithTimeout(ctx, time.Millisecond)
	defer cancel()

	context.AfterFunc(ctx, func() {
		fmt.Println("context done afterFunc task")
	})

	i := 0
	for {
		select {
		case <-ctx.Done():
			fmt.Println("context done")
			return
		default:
			// logic
		}

		i++
		fmt.Println(i, ctx.Value(RequestID("request_id")), ctx.Value("invalid_id"))
		time.Sleep(time.Second)
	}
}
