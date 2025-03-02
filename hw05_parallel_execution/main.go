package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type randNumSlice struct {
	mu     sync.Mutex
	slRand []int
}

func main() {
	//	1 - работа с каналом
	//	2 - сделать структуру с мьютексом
	//	3 - логи в do.once
	//	4 - ассинхронность
	//	5 - select для работы с каналами что бы не было вечного ожидания
	//	6 - обработка ошибок

	sliceMu := &randNumSlice{
		slRand: make([]int, 0),
	}

	readySlice := make([]int, 0)

	chData := make(chan int)
	chErrorSignal := make(chan os.Signal, 1)
	var once sync.Once

	defer close(chData)
	defer close(chErrorSignal)

	wg := sync.WaitGroup{}

	for i := 0; i < 150; i++ {
		wg.Add(1) // <===
		go func() {
			defer wg.Done()
			addRandToSlice(chData, chErrorSignal, &once)
		}()
	}

	wg.Add(1) // <===
	go func() {
		defer wg.Done()
		readySlice = readFromChain(sliceMu, chData, chErrorSignal).slRand
	}()

	wg.Wait()

	println(readySlice)
}

func addRandToSlice(input chan<- int, chErrorSignal chan os.Signal, once *sync.Once) {
	start := mtRand(1, 5)
	end := mtRand(1, 5)

	if start == end {
		errorSet(once, chErrorSignal)
	}

	select {
	case <-chErrorSignal:

	default:
		input <- mtRand(1, 100)
		return
	}
}

func errorSet(once *sync.Once, chErrorSignal chan os.Signal) {
	once.Do(func() {
		signal.Notify(chErrorSignal, syscall.SIGINT, syscall.SIGTERM)
	})
}

func readFromChain(sliceMu *randNumSlice, input <-chan int, chErrorSignal <-chan os.Signal) randNumSlice {

	// Таймер на 4 секунды
	timer := time.NewTimer(4 * time.Second)
	defer timer.Stop()

	for {
		x, ok := <-input
		if !ok {
			break
		}

		select {
		case <-chErrorSignal:
			fmt.Printf("Error")
			break
		case <-timer.C:
			fmt.Println("End.")
			break
		default:
			sliceMu.mu.Lock()
			sliceMu.slRand = append(sliceMu.slRand, x)
			sliceMu.mu.Unlock()
		}
	}

	return randNumSlice{slRand: sliceMu.slRand}
}

func errorPrint() {
	fmt.Println("Error One time")
}

func mtRand(min, max int) int {
	if min > max {
		min, max = max, min
	}
	return rand.Intn(max-min+1) + min
}
