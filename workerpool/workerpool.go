package workerpool

import (
	"fmt"
	"sync"
)

type WorkerManager struct {
	Workers   []*Worker
	TaskQueue chan *Task
	Wg        *sync.WaitGroup
}

type Worker struct {
	Mutex *sync.Mutex
}

type Task struct {
	Name   string
	Params map[string]interface{}
	Exec   func()
}

func Init(numOfWorkers int) *WorkerManager {
	return &WorkerManager{
		Workers:   make([]*Worker, numOfWorkers),
		TaskQueue: make(chan *Task),
		Wg:        &sync.WaitGroup{},
	}
}

func (wm *WorkerManager) AssignTask(tasks ...*Task) {
	for i := range tasks {
		wm.Wg.Add(1)
		wm.TaskQueue <- tasks[i]
	}

	wm.Wg.Wait()
}

func (wm *WorkerManager) Start() {
	for i := range wm.Workers {
		go wm.Workers[i].Run(wm.TaskQueue, wm.Wg)
	}

	fmt.Printf("%d workers are ready !\n", len(wm.Workers))
}

func (w *Worker) Run(taskQueue chan *Task, wg *sync.WaitGroup) {
	for task := range taskQueue {
		task.Exec()
		wg.Done()
	}
}
