package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		return in
	}

	prevChan := in

	for _, stage := range stages {
		stageChan := make(Bi)

		go func(stage Stage, inChan In, outChan Bi) {
			defer close(outChan)

			for {
				select {
				case <-done:
					return
				case v, ok := <-inChan:
					if !ok {
						return
					}

					sOut := stage(makeInChan(v))
					for v2 := range sOut {
						select {
						case outChan <- v2:
						}
					}
				}
			}
		}(stage, prevChan, stageChan)

		prevChan = stageChan
	}

	return prevChan
}

func makeInChan(v interface{}) In {
	ch := make(Bi, 1)
	ch <- v
	close(ch)
	return ch
}
