package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m < 1 {
		m = 1
	}
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	var errCount int
	taskCh := make(chan Task, len(tasks))
	for _, task := range tasks {
		taskCh <- task
	}
	close(taskCh)
	for i := 0; i <= n; i++ {
		wg.Add(1)
		go func(t <-chan Task, wg *sync.WaitGroup) {
			defer wg.Done()
			for task := range t {
				err := task()
				var c int
				if err != nil {
					mu.Lock()
					errCount++
					c = errCount
					mu.Unlock()
				}
				if c >= m {
					return
				}
			}
		}(taskCh, &wg)
	}
	wg.Wait()
	if errCount >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
