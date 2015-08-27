package abutil

import (
	"os"
	"os/signal"
	"syscall"
)

// OnSignal calls the given function on the signals SIGHUP, SIGINT, SIGTERM
// and SIGQUIT
func OnSignal(fn func(os.Signal)) {
	sigc := make(chan os.Signal)

	// Listen for signals
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	// Call the function on each one
	for {
		s := <-sigc
		fn(s)
	}
}
