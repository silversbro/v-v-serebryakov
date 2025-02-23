package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	var once sync.Once

	for i := 0; i < 150; i++ {
		wg.Add(1) // <===
		go func() {
			defer wg.Done()
			if mtRand(0, 150) == i {
				once.Do(errorPrint)
			}
			fmt.Println("go-go-go from:", i)
		}()
	}

	for i := 0; i < 100; i++ {
		wg.Add(1) // <===
		go func() {
			defer wg.Done()
			if mtRand(0, 100) == i {
				once.Do(errorPrint)
			}
			fmt.Println("go-go-go from:", i)
		}()
	}

	wg.Wait()
}

func errorPrint() {
	fmt.Println("Error One time")
}

// mtRand имитирует функцию mt_rand из PHP
func mtRand(min, max int) int {
	if min > max {
		min, max = max, min
	}
	return rand.Intn(max-min+1) + min
}
