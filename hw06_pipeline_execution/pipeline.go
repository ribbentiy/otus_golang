package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	switch len(stages) {
	case 0:
		return in
	case 1:
		f := stages[0]
		return executor(in, done, f)
	default:
		f, s := stages[0], stages[1:]
		return ExecutePipeline(executor(in, done, f), done, s...)
	}
}

func executor(in, done In, fn Stage) Out {
	input := make(Bi)
	res := fn(input)
	go func() {
		defer close(input)
		for i := range in {
			select {
			case <-done:
				return
			default:
				input <- i
			}
		}
	}()
	return res
}
