package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond)
	task(ctx)
}

func task(ctx context.Context) {
	ctx = context.WithoutCancel(ctx)

	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*10)
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
			i++
			fmt.Println(i)
			//time.Sleep(time.Second)
		}

	}
}
