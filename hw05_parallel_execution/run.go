package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	// Place your code here.

	return nil
}

func main() {
	var once sync.Once
	onceBody := func() { fmt.Println("Only once") }
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			once.Do(onceBody)
			wg.Done()
		}()
	}
	wg.Wait()
}
