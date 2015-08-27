package abutil

import (
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestOnSignal(t *testing.T) {
	done := make(chan bool)
	sg := syscall.SIGHUP

	go OnSignal(func(s os.Signal) {
		if s != sg {
			t.Errorf("Expected signal %s, but got %s", sg, s)
		}

		done <- true
	})

	// Send interrupt after 10ms
	time.AfterFunc(10*time.Millisecond, func() {
		syscall.Kill(syscall.Getpid(), sg)
	})
	<-done
}

func TestOnSignalNoChannel(t *testing.T) {
	go OnSignal(func(s os.Signal) {})
}

func OnSignalExamples() {
	OnSignal(func(s os.Signal) {
		fmt.Printf("Got signal %s\n", s)
	})

	// Output: Got signal interrupt
}
