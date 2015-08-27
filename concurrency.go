package abutil

import (
	"fmt"
	"sync"
)

// Parallel runs a given function n times concurrently
// The function get's called with it's number (0-n) for logging purposes
// NOTE: Please set runtime.GOMAXPROCS to runtime.NumCPU()
func Parallel(n int, fn func(int)) {
	var wg sync.WaitGroup
	wg.Add(n)
	defer wg.Wait()

	for i := 0; i < n; i++ {
		go func() {
			fn(n)
			wg.Done()
		}()
	}
}

// The most basic call
func ParallelExample() {
	Parallel(4, func(n int) {
		fmt.Print(n)
	})

	// Output: 0123 in any order
}

// If you need to pass parameters to your function, just wrap it in another
// and call the superior function immeditately.
func ParallelExample_Parameters() {
	fn := func(someParam, someOtherParam string) func(int) {
		return func(n int) {
			fmt.Print(n, someParam, someOtherParam)
		}
	}

	Parallel(4, fn("foo", "bar"))
}
