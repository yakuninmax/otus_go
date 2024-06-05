package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Worker
func worker(jobsQueue <-chan Task, syncNum *int64, wg *sync.WaitGroup) {
	// Update wait group
	defer wg.Done()

	for {
		if task, ok := <-jobsQueue; ok && atomic.LoadInt64(syncNum) >= 0 {
			if task() != nil {
				atomic.AddInt64(syncNum, -1)
			}
		} else {
			return
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {

	// Create wait group
	var wg = sync.WaitGroup{}

	// Check n
	if n <= 0 {
		return nil
	}

	// Check m
	if m < 0 {
		m = 0
	}

	// Create channels and workers
	jobsQueue := make(chan Task, len(tasks))

	var syncNum int64 = int64(m)

	// Start workers
	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(jobsQueue, &syncNum, &wg)
	}

	// Fill work queue
	for _, task := range tasks {
		jobsQueue <- task
	}
	close(jobsQueue)

	wg.Wait()

	if atomic.LoadInt64(&syncNum) >= 0 {
		return nil
	}

	return ErrErrorsLimitExceeded
}
