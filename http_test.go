package abutil

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func mockRequestContext(t *testing.T, fn func(*http.Request)) {
	b := bytes.NewReader([]byte(""))
	r, err := http.NewRequest("GET", "http://some.url", b)
	if err != nil {
		t.Error(err)
	}

	fn(r)
}

func TestRemoteIP(t *testing.T) {
	mockRequestContext(t, func(r *http.Request) {
		ip := "123.456.7.8"

		r.Header = http.Header{
			"X-Real-Ip": []string{ip},
		}

		out := RemoteIP(r)
		if out != ip {
			t.Errorf("Expected %s, but got %s", ip, out)
		}
	})
}

func TestRemoteIPForwardedFor(t *testing.T) {
	mockRequestContext(t, func(r *http.Request) {
		ip := "123.456.7.8"

		r.Header = http.Header{
			"X-Forwarded-For": []string{ip},
		}

		out := RemoteIP(r)
		if out != ip {
			t.Errorf("Expected %s, but got %s", ip, out)
		}
	})
}

func TestRemoteIPRemoteAddr(t *testing.T) {
	mockRequestContext(t, func(r *http.Request) {
		ip := "123.456.7.8"
		r.RemoteAddr = ip

		out := RemoteIP(r)
		if out != ip {
			t.Errorf("Expected %s, but got %s", ip, out)
		}
	})
}

func TestRemoteIPLocalhost(t *testing.T) {
	mockRequestContext(t, func(r *http.Request) {
		ip := "127.0.0.1"
		r.RemoteAddr = "["

		out := RemoteIP(r)
		if out != ip {
			t.Errorf("Expected %s, but got %s", ip, out)
		}
	})
}

func remoteIPMockServe(h http.HandlerFunc) {
	mockRequestContext(nil, func(r *http.Request) {
		r.RemoteAddr = "123.456.7.8"

		w := httptest.NewRecorder()
		h(w, r)
	})
}

func ExampleRemoteIP() {
	someHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("New request from %s\n", RemoteIP(r))
	}

	// Get's called with RemoteAddr = 123.456.7.8
	remoteIPMockServe(someHandler)

	// Output: New request from 123.456.7.8
}

func gracefulServerContext(t *testing.T, fn func(*GracefulServer)) {
	p := 1337
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Foobar"))
	})
	to := 1 * time.Second

	fn(NewGracefulServer(p, h, to))
}

func TestNewGracefulServer(t *testing.T) {
	gracefulServerContext(t, func(s *GracefulServer) {
		if s.server.NoSignalHandling != true {
			t.Error("NoSignalHandling should be true")
		}

		if s.server.Addr != ":1337" {
			t.Error("Didn't set the port correctly")
		}
	})
}

func TestGracefulServerStopped(t *testing.T) {
	gracefulServerContext(t, func(s *GracefulServer) {
		if !s.Stopped() {
			t.Error("Stopped returned false, but shouldn't")
		}

		s.setStopped(false)

		if s.Stopped() {
			t.Error("Stopped returned true, but shouldn't")
		}
	})
}

func TestGracefulServerStop(t *testing.T) {
	done := make(chan bool)

	gracefulServerContext(t, func(s *GracefulServer) {
		time.AfterFunc(10*time.Millisecond, func() {
			s.Stop()

			if !s.Stopped() {
				t.Error("Expected stopped to be true")
			}

			done <- true
		})

		s.ListenAndServe()
	})

	<-done
}

func TestGracefulServerListenAndServe(t *testing.T) {
	done := make(chan bool)

	gracefulServerContext(t, func(s *GracefulServer) {
		time.AfterFunc(10*time.Millisecond, func() {
			if s.Stopped() {
				t.Error("The server should not be stopped after ListenAndServe")
			}

			s.Stop()
			done <- true
		})

		s.ListenAndServe()
	})

	<-done
}

func TestGracefulServerListenAndServeTLS(t *testing.T) {
	done := make(chan bool)

	gracefulServerContext(t, func(s *GracefulServer) {
		time.AfterFunc(10*time.Millisecond, func() {
			if s.Stopped() {
				t.Error("The server should not be stopped after ListenAndServe")
			}

			s.Stop()
			done <- true
		})

		s.ListenAndServeTLS("foo", "bar")
	})

	<-done
}
