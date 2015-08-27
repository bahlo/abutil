package abutil

import (
	"sync"
)

// Parallel runs a given function n times concurrently
// NOTE: Please set runtime.GOMAXPROCS to runtime.NumCPU() for best
// performance
func Parallel(n int, fn func()) {
	var wg sync.WaitGroup
	wg.Add(n)
	defer wg.Wait()

	for i := 0; i < n; i++ {
		go func() {
			fn()
			wg.Done()
		}()
	}
}
