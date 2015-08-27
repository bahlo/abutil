package abutil

import (
	"fmt"
	"os"
	"syscall"
	"testing"
)

func TestOnSignal(t *testing.T) {
	c := make(chan os.Signal)
	done := make(chan bool)

	sg := syscall.SIGINT

	go OnSignal(func(s os.Signal) {
		if s != sg {
			t.Errorf("Expected signal %s, but got %s", sg, s)
		}

		done <- true
	}, c)

	c <- sg
	<-done
}

func OnSignalExamples() {
	OnSignal(func(s os.Signal) {
		fmt.Printf("Got signal %s\n", s)
	})

	// Output: Got signal interrupt
}
