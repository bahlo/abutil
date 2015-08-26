package abutil

import (
	"sync"
)

// Parallel runs a given function n times concurrently
// The function get's called with it's number (0-n) for logging purposes
func Parallel(n int, fn func(int)) {
	var wg sync.WaitGroup

	for i := n; i > 0; i-- {
		wg.Add(1)
		go func() {
			fn(n - i)
			wg.Done()
		}()
	}

	wg.Wait()
}
