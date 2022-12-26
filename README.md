# go-worker-pool
- An implement of go worker-pool : execute multiple tasks (function callings) concurrently
# Example :
```
func main() {
	pool := workerpool.Init(10)
	pool.Start()

	sum := workerpool.NewTask("Sum", func() {
		Sum(10, 20)
	})

	delta := workerpool.NewTask("Delta", func() {
		Delta(20, 10)
	})

	pool.AssignTask(sum, delta)

	fmt.Println(pool.TotalWorkers())
}

func Sum(a, b int) {
	fmt.Println(a + b)
}

func Delta(a, b int) {
	fmt.Println(a - b)
}
```
