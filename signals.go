package abutil

import (
	"os"
	"os/signal"
	"syscall"
)

// OnSignal calls the given function on the signals SIGHUP, SIGINT, SIGTERM
// and SIGQUIT
// You can optionally pass one signal channel to the function and it will use
// this to listen to signals (useful for testing)
func OnSignal(fn func(os.Signal), c ...chan os.Signal) {
	// Get or create chan
	var sigc chan os.Signal
	if len(c) >= 1 {
		sigc = c[0]
	} else {
		sigc = make(chan os.Signal)
	}

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
