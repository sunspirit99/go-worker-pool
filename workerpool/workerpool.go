package workerpool

import (
	"errors"
	"fmt"
	"sync"
)

type WorkerPool struct {
	Workers []*Worker
	Pool    chan *Task
	Wg      *sync.WaitGroup
}

type Worker struct {
	Mutex *sync.Mutex
}

type Task struct {
	Name string
	Exec func()
}

// Initialize a WorkerPool
func Init(numOfWorkers int) *WorkerPool {
	return &WorkerPool{
		Workers: make([]*Worker, numOfWorkers),
		Pool:    make(chan *Task),
		Wg:      &sync.WaitGroup{},
	}
}

// Initialize a task with name and corresponding handler function
func NewTask(name string, handler func()) *Task {
	return &Task{
		Name: name,
		Exec: handler,
	}
}

// Assign a list of tasks for all workers(need to start a pool before)
func (wm *WorkerPool) AssignTask(tasks ...*Task) {
	for i := range tasks {
		wm.Wg.Add(1)
		wm.Pool <- tasks[i]
	}

	wm.Wg.Wait()
}

// Start the pool with n workers
func (wm *WorkerPool) Start() {
	for i := range wm.Workers {
		go wm.Workers[i].run(wm.Pool, wm.Wg)
	}
}

func (w *Worker) run(taskQueue chan *Task, wg *sync.WaitGroup) {
	for task := range taskQueue {
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Recovered. Your task has a panic !")
				}
			}()
			task.Exec()
		}()

		wg.Done()
	}
}

// Add n workers to current pool(if n < 0 then pool will reduced n workers)
func (wm *WorkerPool) AddWorker(n int) error {
	if n >= 0 {
		for i := 1; i <= n; i++ {
			wm.Workers = append(wm.Workers, &Worker{})
		}
	} else {
		lastWorkerId := len(wm.Workers)

		if n > lastWorkerId {
			return errors.New("n must be less than total current number of workers of pool")
		}

		wm.Workers = wm.Workers[:lastWorkerId+n]
	}

	return nil

}

// Get total current number of workers of pool
func (wm *WorkerPool) TotalWorkers() int {
	return len(wm.Workers)
}
