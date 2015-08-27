package abutil

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// BeforeExit calls the given functions once on the signals SIGHUP, SIGINT,
// SIGTERM and SIGQUIT
func BeforeExit(fn func(os.Signal)) {
	sigc := make(chan os.Signal)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	s := <-sigc
	fn(s)
}

func BeforeExitExample() {
	BeforeExit(func(s os.Signal) {
		fmt.Printf("Got signal %s\n", s)
	})

	// Output: Got signal interrupt
}
