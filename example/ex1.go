package example

import (
	"errors"
	"fmt"

	"github.com/sunspirit9999/go-worker-pool/workerpool"
)

var (
	numOfTasks   = 100
	numOfWorkers = 10
)

func New() {
	pool := workerpool.Init(numOfWorkers)
	pool.Start()

	fmt.Println("Total workers : ", pool.TotalWorkers())

	task1 := workerpool.NewTask("Sum", func() interface{} {
		res := Sum(1, 2)
		return res
	})

	var tasks []*workerpool.Task
	for i := 1; i <= numOfTasks; i++ {
		task2 := workerpool.NewTask(fmt.Sprintf("Divide_%d", i), func() interface{} {
			mul, div, err := MultipleAndDivide(1, 1)
			if err != nil {
				return err
			}
			return []int{mul, div}
		})
		tasks = append(tasks, task2)
	}

	pool.AssignTask(append(tasks, task1)...)

	// Get result from the pool
	result := pool.GetResult("Divide_10")
	switch res := result.(type) {
	case error:
		fmt.Printf("type : %T, response : %+v\n", res, res)
	case []int:
		fmt.Printf("type : %T, response : %+v\n", res, res)
	}

	// pool.RefreshWorker(1, 2, 3)

	for _, worker := range pool.Workers {
		fmt.Printf("Worker %d : %+v\n", worker.Id, worker.Result)
	}

	pool.AddWorker(-4)

	for _, worker := range pool.Workers {
		fmt.Printf("Worker %d : %+v\n", worker.Id, worker.Result)
	}

	pool.AddWorker(6)

	for _, worker := range pool.Workers {
		fmt.Printf("Worker %d : %+v\n", worker.Id, worker.Result)
	}

}

func Sum(a, b int) int {
	return a + b
}

func Multiple(a, b int) int {
	return a * b
}

func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("b must be different from 0")
	}
	return a / b, nil
}

func Delta(a, b int) {
	fmt.Println(a - b)
}

func MultipleAndDivide(a, b int) (int, int, error) {
	mul := Multiple(a, b)
	div, err := Divide(a, b)
	if err != nil {
		return mul, 0, err
	}

	return mul, div, nil
}
