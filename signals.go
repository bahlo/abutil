package abutil

import (
	"os"
	"os/signal"
	"syscall"
)

// BeforeExit calls the given functions on various signals
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
