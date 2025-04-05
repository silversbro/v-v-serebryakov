package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

var ErrEmptyTask = errors.New("errors empty task")

type Task func() error

type SafeSlice struct {
	mu    sync.Mutex
	slice []int
}

func (s *SafeSlice) Append(value int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.slice = append(s.slice, value)
}

func (s *SafeSlice) Get() []int {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]int, len(s.slice))
	copy(result, s.slice)

	return result
}

func executeTask(
	tasksCh <-chan Task,
	data *SafeSlice,
	wg *sync.WaitGroup,
	maxError int64,
	errorCount *int64,
) {
	defer wg.Done()

	for task := range tasksCh {
		err := task()
		if err != nil {
			if atomic.AddInt64(errorCount, 1) >= maxError {
				return
			}
		} else {
			data.Append(1)
		}
	}
}

func Run(tasks []Task, n, m int) error {
	if len(tasks) == 0 {
		return ErrEmptyTask
	}

	if m == 0 {
		return ErrErrorsLimitExceeded
	}

	jobs := make(chan Task, len(tasks))
	var err error
	var wg sync.WaitGroup
	safeSlice := &SafeSlice{}
	maxError := int64(m)
	var errorCount int64

	for i := 0; i < n; i++ {
		wg.Add(1)
		go executeTask(jobs, safeSlice, &wg, maxError, &errorCount)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, task := range tasks {
			if atomic.LoadInt64(&errorCount) > maxError {
				return
			}

			jobs <- task
		}
		close(jobs)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			loadCountError := atomic.LoadInt64(&errorCount)
			if loadCountError >= maxError {
				err = ErrErrorsLimitExceeded

				return
			}

			if (int64(len(safeSlice.Get())) + loadCountError) == int64(len(tasks)) {
				return
			}
		}
	}()

	wg.Wait()

	fmt.Printf("Successful tasks: %d\n", len(safeSlice.Get()))
	fmt.Printf("Errors: %d\n", errorCount)

	return err
}
