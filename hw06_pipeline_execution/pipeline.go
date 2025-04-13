package hw06pipelineexecution

import "sync"

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// Если нет стейджей, просто возвращаем входной канал
	if len(stages) == 0 {
		return in
	}

	// Создаем выходной канал
	out := make(chan interface{})

	var wg sync.WaitGroup

	// Запускаем обработку в отдельной горутине
	go func() {
		defer close(out)

		// Создаем каналы для каждого стейджа
		stageChans := make([]chan interface{}, len(stages)+1)
		for i := range stageChans {
			stageChans[i] = make(chan interface{})
		}

		// Запускаем все стейджи в отдельных горутинах
		for i, stage := range stages {
			wg.Add(1)
			go func(index int, s Stage) {
				defer wg.Done()
				defer close(stageChans[index+1])

				// Запускаем стейдж и перенаправляем его вывод
				for val := range s(stageChans[index]) {
					select {
					case <-done:
						return
					case stageChans[index+1] <- val:
					}
				}
			}(i, stage)
		}

		// Перенаправляем входные данные в первый стейдж
		go func() {
			defer close(stageChans[0])
			for {
				select {
				case <-done:
					return
				case val, ok := <-in:
					if !ok {
						return
					}
					select {
					case <-done:
						return
					case stageChans[0] <- val:
					}
				}
			}
		}()

		// Перенаправляем вывод последнего стейджа в выходной канал
		for val := range stageChans[len(stageChans)-1] {
			select {
			case <-done:
				return
			case out <- val:
			}
		}

		// Ждем завершения всех стейджей
		wg.Wait()
	}()

	return out
}
