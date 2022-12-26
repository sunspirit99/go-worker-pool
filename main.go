package main

import (
	"fmt"

	wm "github.com/sunspirit9999/go-worker-manager/workerpool"
)

func main() {
	boss := wm.Init(10)
	boss.Start()

	var tasks = []*wm.Task{
		{
			Name: "Sum",
			Exec: func() {
				Sum(10, 20)
			},
		},
		{
			Name: "Delta",
			Exec: func() {
				Delta(50, 20)
			},
		},
	}

	boss.AssignTask(tasks...)
}

func Sum(a, b int) {
	fmt.Println(a + b)
}

func Delta(a, b int) {
	fmt.Println(a - b)
}
