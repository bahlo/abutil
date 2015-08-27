package abutil

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestParallel(t *testing.T) {
	var m sync.Mutex
	var wg sync.WaitGroup

	counter := 0

	wg.Add(2)
	go Parallel(2, func(n int) {
		m.Lock()
		counter++
		m.Unlock()
		wg.Done()
	})

	wg.Wait()

	if counter != 2 {
		t.Errorf("Expected counter to be %d, but got %d", 2, counter)
	}
}

func TestParallelCounter(t *testing.T) {
	var m sync.Mutex
	var wg sync.WaitGroup

	sum := 0
	wg.Add(4)
	go Parallel(4, func(n int) {
		m.Lock()
		sum += n
		m.Unlock()
		wg.Done()
	})

	wg.Wait()

	if sum != 6 {
		t.Errorf("Expected sum to be %d, but got %d", 6, sum)
	}
}

func TestParallelTiming(t *testing.T) {
	var m sync.Mutex
	counter := 0

	go Parallel(4, func(n int) {
		time.Sleep(time.Duration(10*n) * time.Millisecond)
		m.Lock()
		counter++
		m.Unlock()
	})

	done := make(chan bool)
	time.AfterFunc(25*time.Millisecond, func() {
		m.Lock()
		if counter != 2 {
			t.Errorf("Expected counter to be %d, but got %d", 2, counter)
		}
		m.Unlock()
		done <- true
	})

	<-done
}

// The most basic call
func ParallelExamples() {
	Parallel(4, func(n int) {
		fmt.Print(n)
	})

	// Output: 0123 in any order
}

// If you need to pass parameters to your function, just wrap it in another
// and call the superior function immeditately.
func ParallelExamples_Parameters() {
	fn := func(someParam, someOtherParam string) func(int) {
		return func(n int) {
			fmt.Print(n, someParam, someOtherParam)
		}
	}

	Parallel(4, fn("foo", "bar"))
}
