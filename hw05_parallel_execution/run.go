package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

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
	// Возвращаем копию слайса, чтобы избежать race condition
	result := make([]int, len(s.slice))
	copy(result, s.slice)

	return result
}

func executeTask(
	tasksCh <-chan Task,
	data *SafeSlice,
	errCh chan<- error,
	wg *sync.WaitGroup,
	maxError int,
	errorCount *int,
	errorCountMutex *sync.Mutex,
	stopSignal *bool,
	stopSignalMutex *sync.Mutex,
) {
	defer wg.Done()

	for task := range tasksCh {
		// Проверяем, не было ли сигнала остановки
		stopSignalMutex.Lock()
		if *stopSignal {
			stopSignalMutex.Unlock()

			return
		}
		stopSignalMutex.Unlock()

		err := task()
		if err != nil {
			errorCountMutex.Lock()
			*errorCount++ // Увеличиваем счетчик ошибок
			errorCountMutex.Unlock()
			errCh <- err

			errorCountMutex.Lock()
			if *errorCount > maxError {
				stopSignalMutex.Lock()
				*stopSignal = true // Устанавливаем сигнал остановки
				stopSignalMutex.Unlock()
				errorCountMutex.Unlock()
				return
			}
			errorCountMutex.Unlock()
		} else {
			data.Append(1)
		}
	}
}

func Run(tasks []Task, n, m int) error {
	errCh := make(chan error, len(tasks))
	defer close(errCh)

	jobs := make(chan Task, len(tasks))
	var err error

	var wg sync.WaitGroup
	safeSlice := &SafeSlice{}

	var errorCount int
	var errorCountMutex sync.Mutex
	var stopSignal bool
	var stopSignalMutex sync.Mutex

	for i := 0; i < n; i++ {
		wg.Add(1)
		go executeTask(jobs, safeSlice, errCh, &wg, m, &errorCount, &errorCountMutex, &stopSignal, &stopSignalMutex)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, task := range tasks {
			jobs <- task
		}
		close(jobs)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		timer := time.NewTimer(10 * time.Second)
		defer timer.Stop()

		for {
			select {
			case <-timer.C:
				stopSignalMutex.Lock()
				stopSignal = true
				stopSignalMutex.Unlock()
				err = errors.New("timeout")

				return
			case <-errCh:
				errorCountMutex.Lock()
				if errorCount > m {
					stopSignalMutex.Lock()
					stopSignal = true
					stopSignalMutex.Unlock()
					err = ErrErrorsLimitExceeded
					errorCountMutex.Unlock()

					return
				}

				errorCountMutex.Unlock()
			default:
				stopSignalMutex.Lock()
				if stopSignal {
					stopSignalMutex.Unlock()

					return
				}

				stopSignalMutex.Unlock()

				if (len(safeSlice.Get()) + errorCount) == len(tasks) {

					return
				}

				// Небольшая задержка, чтобы не загружать процессор
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	wg.Wait()

	fmt.Printf("Successful tasks: %d\n", len(safeSlice.Get()))
	fmt.Printf("Errors: %d\n", errorCount)

	return err
}
