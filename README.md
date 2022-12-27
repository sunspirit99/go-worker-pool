# go-worker-pool
- An implement of go worker-pool : execute multiple tasks (function callings) concurrently
# Example :
```
func main() {
	// Init a pool with 10 workers
	pool := workerpool.Initialize(10)
	
	// Start the pool
	pool.Start()
	
	// Get total workers of pool
	fmt.Println(pool.TotalWorkers())

	// Init task : sum of 2 number
	sum := workerpool.NewTask("Sum", func() {
		Sum(10, 20)
	})
	
	// Init task : delta of 2 number
	delta := workerpool.NewTask("Delta", func() {
		Delta(20, 10)
	})
	
	// Assign the tasks to the pool for the workers inside to take out and process it
	pool.AssignTask(context.Background(), sum, delta)
}


func Sum(a, b int) {
	fmt.Println(a + b)
}

func Delta(a, b int) {
	fmt.Println(a - b)
}
```
