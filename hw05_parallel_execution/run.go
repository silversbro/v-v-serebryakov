package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Структура для хранения данных с мьютексом
type SafeSlice struct {
	mu    sync.Mutex
	slice []int
}

// Добавление значения в слайс
func (s *SafeSlice) Append(value int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.slice = append(s.slice, value)
}

// Получение слайса
func (s *SafeSlice) Get() []int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.slice
}

func executeTask(
	tasksCh <-chan Task,
	dataCh chan int,
	errCh chan int,
	wg *sync.WaitGroup,
	stopChan chan error,
	maxErrors int,
) {
	defer wg.Done()
	fmt.Println("Start task execute ")

	for {
		select {
		case task, ok := <-tasksCh:
			if !ok {
				return
			}

			err := task()

			if err != nil {
				errCh <- 1
				fmt.Printf("Error in task: %v\n", err)

				if len(errCh) > maxErrors {
					stopChan <- ErrErrorsLimitExceeded

					return
				}

				return
			}

			dataCh <- 1

			return
		case <-stopChan:
			fmt.Println("Error in task stop chain")

			return
		case <-errCh:
			fmt.Println("Error in task")

			if len(errCh) > maxErrors {
				stopChan <- ErrErrorsLimitExceeded

				return
			}
		}
	}
}

func Run(tasks []Task, n, m int) error {
	dataCh := make(chan int)
	errCh := make(chan int)
	var err error

	jobs := make(chan Task)
	stopChan := make(chan error)
	defer close(errCh)
	defer close(stopChan)
	defer close(dataCh)

	var wg sync.WaitGroup
	safeSlice := SafeSlice{}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go executeTask(jobs, dataCh, errCh, &wg, stopChan, m)
	}

	for _, task := range tasks {
		jobs <- task
	}
	defer close(jobs)

	// Горутина для чтения данных из канала
	wg.Add(1)
	go func() {
		defer wg.Done()

		timer := time.NewTimer(30 * time.Second)
		defer timer.Stop()

		for {
			select {
			case num, ok := <-dataCh:
				if !ok {
					return
				}

				safeSlice.Append(num)
			case <-timer.C:
				// Время вышло
				total := len(safeSlice.Get()) + len(errCh)
				fmt.Printf("Всего задач и ошибок: %v\n", total)

				panic("Timeout.")
			case err = <-stopChan:
				fmt.Printf("Error: %v\n", err)

				return
			}
		}
	}()

	// Ожидание завершения всех горутин
	wg.Wait()
	total := len(safeSlice.Get()) + len(errCh)

	fmt.Printf("Выполнилось: %v\n", total)

	return err
}
