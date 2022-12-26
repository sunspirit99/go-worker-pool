package example

import (
	"fmt"
	"math/rand"

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

	var task []*workerpool.Task

	for i := 1; i <= numOfTasks; i++ {
		task = append(task, workerpool.NewTask("Sum", func() {
			a := rand.Intn(10)
			b := rand.Intn(10)
			Sum(a, b)
		}))
	}

	pool.AssignTask(task...)

}

func Sum(a, b int) {
	fmt.Println(a + b)
}

func Delta(a, b int) {
	fmt.Println(a - b)
}
