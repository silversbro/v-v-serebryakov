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

// Функция для генерации случайного числа и возможной ошибки
func executeTask(
	tasks <-chan Task,
	ch chan<- int,
	errCh chan int,
	wg *sync.WaitGroup,
	stopChan chan error,
	maxErrors int,
) {
	defer wg.Done()

	for {
		select {
		case <-stopChan:
			return
		case <-errCh:
			return
		case task := <-tasks:
			err := task()

			if err != nil {
				errCh <- 1
				return
			}

			select {
			case ch <- 1:
			case <-errCh:
				if len(errCh) > maxErrors {
					stopChan <- ErrErrorsLimitExceeded

					return
				}

				continue
			case <-stopChan:
				return
			}
		default:
			return
		}
	}
}

func Run(tasks []Task, n, m int) error {
	dataCh := make(chan int)
	errCh := make(chan int)

	jobs := make(chan Task, len(tasks))
	stopChan := make(chan error)
	defer close(dataCh)
	defer close(errCh)

	// WaitGroup для ожидания завершения всех горутин
	var wg sync.WaitGroup

	// Структура для хранения данных
	safeSlice := SafeSlice{}

	// Запуск 100 горутин для генерации чисел
	for i := 0; i < n; i++ {
		wg.Add(1)
		go executeTask(jobs, dataCh, errCh, &wg, stopChan, m)
	}

	for _, task := range tasks {
		jobs <- task
	}

	close(jobs)

	// Горутина для чтения данных из канала
	wg.Add(1)
	go func() {
		defer wg.Done()

		timer := time.NewTimer(2 * time.Second)
		defer timer.Stop()

		for {
			select {
			case num := <-dataCh:
				safeSlice.Append(num)
			case <-timer.C:
				// Время вышло
				fmt.Println("Timeout.")
				close(stopChan)
				return
			case err = <-stopChan:
				close(stopChan)
				fmt.Printf("Error: %v\n", err)

				return err
			}
		}
	}()

	// Ожидание завершения всех горутин
	wg.Wait()

	fmt.Printf("Выполнилось: %v\n", len(safeSlice.Get()))

	return error
}
