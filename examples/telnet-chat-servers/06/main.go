package main

import (
	"context"
	"fmt"
	"multictx/MultivalContext"
	"time"
)

func main() {

	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, time.Second*10)
	ctx = context.WithValue(ctx, "name", "golang")
	f1(ctx)

}

func f1(ctx context.Context) {
	kv := map[any]any{
		"Учебный центр": "ОТУС",
		"учим язык":     "GO",
	}

	newCtx := MultivalContext.WithMultivalContext(ctx, kv)

	f2(newCtx)

}

func f2(ctx context.Context) {
	fmt.Printf("Получен контекст:")
	select {
	case <-ctx.Done():
		fmt.Println("Контекст отменен")
	default:
		fmt.Println("Контекст актуален")
	}

	if t, ok := ctx.Deadline(); ok {
		fmt.Printf("У контекста есть дедлайн, %v \n", t)
	} else {
		fmt.Printf("Контекста у дедлайна нет \n")
	}

	if v := ctx.Value("name"); v != nil {
		fmt.Printf("Ключ %v в контексте имеет значение %v \n", "name", v)
	} else {
		fmt.Printf("Ключ %v в контексте отсутствует \n", "name")
	}

	if v := ctx.Value("Учебный центр"); v != nil {
		fmt.Printf("Ключ %v в контексте имеет значение %v \n", "Учебный центр", v)
	} else {
		fmt.Printf("Ключ %v в контексте отсутствует \n", "Учебный центр")
	}

	if v := ctx.Value("ключ"); v != nil {
		fmt.Printf("Ключ %v в контексте имеет значение %v \n", "ключ", v)
	} else {
		fmt.Printf("Ключ %v в контексте отсутствует \n", "ключ")
	}

}

func f3(ctx context.Context) {

}
