package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Worker function.
func worker(taskQueue <-chan Task, errorsCount *int64, waitGroup *sync.WaitGroup) error {
	// Done wait group.
	defer waitGroup.Done()

	// Run tasks.
	for {
		// Get task, and check errorsCount.
		// Run task if errorsCount >= 0.
		if task, ok := <-taskQueue; ok && atomic.LoadInt64(errorsCount) >= 0 {
			// Check task result. If error, decrease errorsCounter.
			if task() != nil {
				atomic.AddInt64(errorsCount, -1)
			}
		} else {
			return ErrErrorsLimitExceeded
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	// Check n & m values.
	if n <= 0 {
		return errors.New("n must be greater than 0")
	}

	if m < 0 {
		m = 0
	}

	// Create channel.
	taskQueue := make(chan Task, len(tasks))

	// Create wait group.
	waitGroup := sync.WaitGroup{}

	// Set errors counter.
	errorsCount := int64(m)

	// Run n workers.
	for i := 0; i < n; i++ {
		waitGroup.Add(1)
		go worker(taskQueue, &errorsCount, &waitGroup)
	}

	// Send tasks to channel.
	for _, task := range tasks {
		taskQueue <- task
	}
	close(taskQueue)

	// Wait for jobs done.
	waitGroup.Wait()

	// Check errorsCount.
	if atomic.LoadInt64(&errorsCount) < 0 {
		return ErrErrorsLimitExceeded
	}

	fmt.Println(errorsCount)

	return nil
}
