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

// init a WorkerPool
func Init(numOfWorkers int) *WorkerPool {
	return &WorkerPool{
		Workers: make([]*Worker, numOfWorkers),
		Pool:    make(chan *Task),
		Wg:      &sync.WaitGroup{},
	}
}

// init a task with name and execution func
func NewTask(name string, exec func()) *Task {
	return &Task{
		Name: name,
		Exec: exec,
	}
}

// assign a list of tasks for all workers (need to start WorkerPool before)
func (wm *WorkerPool) AssignTask(tasks ...*Task) {
	for i := range tasks {
		wm.Wg.Add(1)
		wm.Pool <- tasks[i]
	}

	wm.Wg.Wait()
}

// start WorkerPool
func (wm *WorkerPool) Start() {
	for i := range wm.Workers {
		go wm.Workers[i].run(wm.Pool, wm.Wg)
	}

	fmt.Printf("%d workers are ready !\n", len(wm.Workers))
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

// Add n workers to current pool
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
