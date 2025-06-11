package hw06pipelineexecution

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

	current := in
	for _, stage := range stages {
		// Запускаем каждый стейдж, передавая ему текущий канал
		current = runStage(stage, current, done)
	}

	return current
}

func runStage(stage Stage, in In, done In) Out {
	out := make(Bi)

	go func() {
		stageOut := stage(in)
		defer close(out)

		for {
			select {
			case <-done:
				go drainChannel(stageOut)

				return
			case val, ok := <-stageOut:
				if !ok {
					return
				}
				select {
				case out <- val:
					continue
				case <-done:
					go drainChannel(stageOut)

					return
				}
			}
		}
	}()

	return out
}

func drainChannel(ch <-chan interface{}) {
	//nolint
	for range ch {
	}
}
