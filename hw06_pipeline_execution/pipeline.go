package hw06pipelineexecution

import "sync"

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

	var wg sync.WaitGroup
	prevChan := in

	for _, stage := range stages {
		stageChan := make(Bi)
		wg.Add(1)

		go func(stage Stage, inChan In, outChan Bi) {
			defer wg.Done()
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
						case <-done:
							return
						case outChan <- v2:
						}
					}
				}
			}
		}(stage, prevChan, stageChan)

		prevChan = stageChan
	}

	out := make(Bi)
	go func() {
		defer close(out)
		for v := range prevChan {
			select {
			case <-done:
				return
			case out <- v:
			}
		}
	}()

	go func() {
		wg.Wait()
	}()

	return out
}

func makeInChan(v interface{}) In {
	ch := make(Bi, 1)
	ch <- v
	close(ch)
	return ch
}
