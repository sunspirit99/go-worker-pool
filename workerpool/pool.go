package workerpool

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type WorkerPool struct {
	Workers []*Worker
	Pool    chan *Task
	Wg      *sync.WaitGroup
}

type Worker struct {
	Id     int
	Result map[string]interface{}
}

type Task struct {
	Name    string
	Execute func() interface{}
}

// Initialize a WorkerPool
func Init(numOfWorkers int) *WorkerPool {
	var workers []*Worker
	for i := 1; i <= numOfWorkers; i++ {
		workers = append(workers, &Worker{
			Id: i,
		})
	}
	return &WorkerPool{
		Workers: workers,
		Pool:    make(chan *Task),
		Wg:      &sync.WaitGroup{},
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
	w.Result = make(map[string]interface{})

	for task := range taskQueue {
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Recovered. Your task has a panic !")
				}
			}()

			result := task.Execute()

			if result != nil {
				if reflect.ValueOf(result).IsValid() {
					w.Result[task.Name] = result
				}
			}
		}()

		wg.Done()
	}
}

func (wm *WorkerPool) GetResult(taskName string) interface{} {
	for i := range wm.Workers {
		if _, ok := wm.Workers[i].Result[taskName]; ok {
			return wm.Workers[i].Result[taskName]
		}
	}
	return nil
}

// Add n workers to current pool(if n < 0 then pool will remove n last workers)
func (wm *WorkerPool) AddWorker(n int) error {
	lastWorkerId := len(wm.Workers)
	if n >= 0 {
		id := lastWorkerId
		for i := 1; i <= n; i++ {
			id += 1
			wm.Workers = append(wm.Workers, &Worker{
				Id:     id,
				Result: make(map[string]interface{}),
			})
		}
	} else {
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

// worker refreshing will free result stored at a list of workers (free all workers if calling function without params)
func (wm *WorkerPool) RefreshWorker(id ...int) {
	if len(id) == 0 {
		for i := range wm.Workers {
			wm.Workers[i].Result = nil
		}
	} else {
		for i := 0; i <= len(id); i++ {
			wm.Workers[i].Result = nil
		}
	}
}
