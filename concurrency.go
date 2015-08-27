package abutil

import (
	"sync"
)

// Parallel runs a given function n times concurrently
// The function get's called with it's number (0-n) for logging purposes
// NOTE: Please set runtime.GOMAXPROCS to runtime.NumCPU()
func Parallel(n int, fn func(int)) {
	var wg sync.WaitGroup
	wg.Add(n)
	defer wg.Wait()

	var m sync.Mutex
	counter := 0

	for i := 0; i < n; i++ {
		go func() {
			// Lock mutex to prevent multiple functions with the same index
			m.Lock()

			// Run function
			fn(counter)

			// Increment counter
			counter++

			m.Unlock()
			wg.Done()
		}()
	}
}
