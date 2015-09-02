package abutil

import (
	"fmt"
	"sync"
	"testing"
)

func TestParallel(t *testing.T) {
	var wg sync.WaitGroup
	var m sync.Mutex

	counter := 0

	wg.Add(2)
	go Parallel(2, func() {
		m.Lock()
		defer func() {
			m.Unlock()
			wg.Done()
		}()

		counter++
	})
	wg.Wait()

	if counter != 2 {
		t.Errorf("Expected counter to be %d, but got %d", 2, counter)
	}
}

// The most basic call
func ExampleParallel() {
	var m sync.Mutex

	c := 0
	Parallel(4, func() {
		m.Lock()
		defer m.Unlock()

		fmt.Print(c)
		c++
	})

	// Output: 0123
}

// If you need to pass parameters to your function, just wrap it in another
// and call the superior function immeditately.
func ExampleParallel_parameters() {
	fn := func(someParam, someOtherParam string) func() {
		return func() {
			fmt.Print(someParam, someOtherParam)
		}
	}

	Parallel(4, fn("foo", "bar"))
}
