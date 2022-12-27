package example

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"time"

	"github.com/sunspirit9999/go-worker-pool/workerpool"
)

var (
	numOfTasks   = 100
	numOfWorkers = 10
)

func New() {

	pool := workerpool.Initialize(numOfWorkers)
	pool.Start()

	fmt.Println("Total workers : ", pool.GetTotalWorkers())

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute) // set timeout
	defer cancel()

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

	err := pool.AssignTask(ctx, tasks...)
	if err != nil {
		fmt.Println(err)
	}

	// Get result from the pool
	result := pool.GetResult("Divide_10")
	switch res := result.(type) {
	case error:
		fmt.Printf("type : %T, response : %+v\n", res, res)
	case []int:
		fmt.Printf("type : %T, response : %+v\n", res, res)
	}

	pool.Close()
	fmt.Println("total worker's threads :", runtime.NumGoroutine()-1) // ignore main thread

	pool.AddWorker(5)
	fmt.Println("total worker's threads :", runtime.NumGoroutine()-1) // ignore main thread

	pool.AssignTask(ctx, tasks...)

	// Get result from the pool
	result = pool.GetResult("Divide_10")
	switch res := result.(type) {
	case error:
		fmt.Printf("type : %T, response : %+v\n", res, res)
	case []int:
		fmt.Printf("type : %T, response : %+v\n", res, res)
	}

	fmt.Println("total worker's threads :", runtime.NumGoroutine()-1) // ignore main thread

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
