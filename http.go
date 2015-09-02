package abutil

import (
	"crypto/tls"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"

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

// GracefulServer is basically graceful.Server (github.com/tylerb/graceful),
// but adds a state variable to check if stopped and doesn't listen on
// signals (use OnSignal instead)
type GracefulServer struct {
	*graceful.Server

	// stopped determines if the server is stopped
	stopped bool

	// locker controls the access to running
	locker sync.Locker
}

// NewGracefulServer creates a new GracefulServer with the given handler,
// which listens on the given port.
func NewGracefulServer(p int, h http.Handler) *GracefulServer {
	var m sync.Mutex
	s := &GracefulServer{
		Server: &graceful.Server{
			Server: &http.Server{
				Addr:    ":" + strconv.Itoa(p),
				Handler: h,
			},
			NoSignalHandling: true,
		},
		stopped: true,
		locker:  &m,
	}

	s.Server.ShutdownInitiated = func() { s.setStopped(true) }

	return s
}

// Stopped returns if the server is stopped
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

// Serve is equivalent to http.Server.Serve with graceful shutdown enabled
func (g *GracefulServer) Serve(l net.Listener) error {
	g.setStopped(false)
	return g.Server.Serve(l)
}

// ListenAndServe is equivalent to http.Server.ListenAndServe with graceful
// shutdown enabled
func (g *GracefulServer) ListenAndServe() error {
	g.setStopped(false)
	return g.Server.ListenAndServe()
}

// ListenAndServeTLS is equivalent to http.Server.ListenAndServeTLS with
// graceful shutdown enabled
func (g *GracefulServer) ListenAndServeTLS(cf, kf string) error {
	g.setStopped(false)
	return g.Server.ListenAndServeTLS(cf, kf)
}

// ListenAndServeTLSConfig is equivalent to
// http.Server.ListenAndServeTLSConfig with graceful shutdown enabled
func (g *GracefulServer) ListenAndServeTLSConfig(c *tls.Config) error {
	g.setStopped(false)
	return g.Server.ListenAndServeTLSConfig(c)
}
