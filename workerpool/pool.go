package workerpool

import (
	"context"
	"errors"
	"log"
	"sync"
)

type WorkerPool struct {
	TotalWorkers int
	Pool         chan *Task
	Wg           *sync.WaitGroup
	Mutex        *sync.Mutex
	Result       map[string]interface{}
	Status       int // 1 is active, 0 is waiting, -1 is inactive
	closeSignal  chan struct{}
	done         chan struct{}
}

type Task struct {
	Name    string
	Execute func() interface{}
}

// Initialize a WorkerPool
func Initialize(numOfWorkers int) *WorkerPool {

	return &WorkerPool{
		TotalWorkers: numOfWorkers,
		Pool:         make(chan *Task),
		Result:       make(map[string]interface{}),
		Wg:           &sync.WaitGroup{},
		Mutex:        &sync.Mutex{},
		Status:       -1,
		closeSignal:  make(chan struct{}, numOfWorkers),
		done:         make(chan struct{}, numOfWorkers),
	}
}

// Initialize a task with name and corresponding handler function
func NewTask(name string, handler func() interface{}) *Task {
	return &Task{
		Name:    name,
		Execute: handler,
	}
}

// Assign a list of tasks for all workers(need to start a pool before)
func (wp *WorkerPool) AssignTask(ctx context.Context, tasks ...*Task) error {

	wp.Status = 1

	for i := range tasks {
		wp.Wg.Add(1)

		wp.Pool <- tasks[i]
	}

	wp.Wg.Wait()

	return nil
}

// Free all workers in the pool
func (wp *WorkerPool) Close() {

	for i := 1; i <= wp.TotalWorkers; i++ {
		wp.closeSignal <- struct{}{}
	}

	var totalClosedWorkers int

	for {
		if totalClosedWorkers == wp.TotalWorkers {
			wp.Status = -1
			wp.TotalWorkers = 0
			wp.Wg = &sync.WaitGroup{}
			return
		}
		<-wp.done
		totalClosedWorkers += 1
	}

}

// Start the pool with n workers
func (wp *WorkerPool) Start() {
	if wp.Pool == nil {
		return
	}

	for i := 1; i <= wp.TotalWorkers; i++ {
		go wp.run()
	}

	wp.Status = 0
}

// Get result of executed tasks by workers
func (wp *WorkerPool) GetResult(taskName string) interface{} {
	if _, ok := wp.Result[taskName]; ok {
		return wp.Result[taskName]
	}

	return errors.New("task's result not found")
}

// Add n workers to current pool(if n < 0 then pool will reduced n workers)
func (wp *WorkerPool) AddWorker(n int) error {
	if wp.TotalWorkers < -n {
		return errors.New("if n < 0, n must be less than total current workers")
	}
	wp.TotalWorkers += n
	wp.Start()
	return nil
}

// Get total current number of workers of pool
func (wp *WorkerPool) GetTotalWorkers() int {
	return wp.TotalWorkers
}

// Check current status of pool
func (wp *WorkerPool) CheckStatus() string {
	switch wp.Status {
	case 1:
		return "active"
	case 0:
		return "waiting"
	default:
		return "inactive"
	}
}

// Remove all current results
func (wp *WorkerPool) RefreshResult() {
	wp.Result = nil
}

func (wp *WorkerPool) run() {

	for {
		select {
		case <-wp.closeSignal:
			defer func() {
				wp.done <- struct{}{}
			}()

			return

		case task := <-wp.Pool:
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Recovered. Your task %s has a panic !\n", task.Name)
					}
				}()

				result := task.Execute()

				wp.Mutex.Lock()
				if result != nil {
					wp.Result[task.Name] = result
				}
				wp.Mutex.Unlock()
			}()

			wp.Wg.Done()
		}
	}

}
