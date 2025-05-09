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

	out := make(chan interface{})
	var wg sync.WaitGroup

	stageChans := make([]chan interface{}, len(stages)+1)
	for i := range stageChans {
		stageChans[i] = make(chan interface{})
	}

	for i, stage := range stages {
		wg.Add(1)
		go func(index int, s Stage) {
			defer wg.Done()
			defer close(stageChans[index+1])

			for {
				select {
				case <-done:
					return
				case v, ok := <-stageChans[index]:
					if !ok {
						return
					}
					sOut := s(makeInChan(v))
					for v2 := range sOut {
						select {
						case <-done:
							return
						case stageChans[index+1] <- v2:
						}
					}
				}
			}
		}(i, stage)
	}

	go func() {
		defer close(stageChans[0])
		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case stageChans[0] <- v:
				}
			}
		}
	}()

	go func() {
		defer close(out)
		for v := range stageChans[len(stageChans)-1] {
			select {
			case <-done:
				return
			case out <- v:
			}
		}
	}()

	return out
}

func makeInChan(v interface{}) In {
	ch := make(chan interface{}, 1)
	ch <- v
	close(ch)
	return ch
}
