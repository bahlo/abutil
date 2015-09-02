package abutil

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tylerb/graceful"
)

// RemoteIP tries to get the remote ip and returns it or ""
func RemoteIP(r *http.Request) string {
	a := r.Header.Get("X-Real-IP")

	if a == "" {
		a = r.Header.Get("X-Forwarded-For")
	}

	if a == "" {
		a = strings.SplitN(r.RemoteAddr, ":", 2)[0]

		// Check localhost
		if a == "[" {
			a = "127.0.0.1"
		}
	}

	return a
}

// GracefulServer is a wrapper around graceful.Server from
// github.com/tylerb/graceful, but adds a running variable and a mutex for
// controlling the access to it.
type GracefulServer struct {
	server *graceful.Server

	// stopped determines if the server is stopped
	stopped bool

	// locker controls the access to running
	locker sync.Locker
}

// NewGracefulServer creates a new GracefulServer with the given handler,
// which listens on the given port. When Stop() ist called, it waits until
// the timeout is finished or all connections are closed (whatever comes first)
func NewGracefulServer(p int, h http.Handler, t time.Duration) *GracefulServer {
	var m sync.Mutex
	return &GracefulServer{
		server: &graceful.Server{
			Server: &http.Server{
				Addr:    ":" + strconv.Itoa(p),
				Handler: h,
			},
			NoSignalHandling: true,
			Timeout:          t,
		},
		stopped: true,
		locker:  &m,
	}
}

// Stopped returns if the server is running
func (g *GracefulServer) Stopped() bool {
	g.locker.Lock()
	defer g.locker.Unlock()

	return g.stopped
}

func (g *GracefulServer) setStopped(r bool) {
	g.locker.Lock()
	g.stopped = r
	g.locker.Unlock()
}

// Stop stops the server (this may last up to the server timeout)
func (g *GracefulServer) Stop() {
	g.setStopped(true)

	g.server.Stop(g.server.Timeout)
}

// Serve is equivalent to http.Server.Serve, but the server is stoppable
func (g *GracefulServer) Serve(l net.Listener) error {
	g.setStopped(false)
	return g.server.Serve(l)
}

// ListenAndServe is equivalent to http.Server.Serve, but the server is
// stoppable
func (g *GracefulServer) ListenAndServe() error {
	g.setStopped(false)
	return g.server.ListenAndServe()
}

// ListenAndServeTLS is equivalent to http.Server.Serve, but the server is
// stoppable
func (g *GracefulServer) ListenAndServeTLS(cf, kf string) error {
	g.setStopped(false)
	return g.server.ListenAndServeTLS(cf, kf)
}
