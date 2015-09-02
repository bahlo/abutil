package abutil

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net"
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
	ip := "123.456.7.8"

	mockRequestContext(t, func(r *http.Request) {
		r.Header = http.Header{
			"X-Real-Ip": []string{ip},
		}

		out := RemoteIP(r)
		if out != ip {
			t.Errorf("Expected %s, but got %s", ip, out)
		}
	})

	mockRequestContext(t, func(r *http.Request) {
		r.Header = http.Header{
			"X-Forwarded-For": []string{ip},
		}

		out := RemoteIP(r)
		if out != ip {
			t.Errorf("Expected %s, but got %s", ip, out)
		}
	})

	mockRequestContext(t, func(r *http.Request) {
		r.RemoteAddr = ip

		out := RemoteIP(r)
		if out != ip {
			t.Errorf("Expected %s, but got %s", ip, out)
		}
	})

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

	fn(NewGracefulServer(p, h))
}

func TestGracefulServer(t *testing.T) {
	gracefulServerContext(t, func(s *GracefulServer) {
		if s.Server.NoSignalHandling != true {
			t.Error("NoSignalHandling should be true")
		}

		if s.Server.Addr != ":1337" {
			t.Error("Didn't set the port correctly")
		}
	})

	gracefulServerContext(t, func(s *GracefulServer) {
		if !s.Stopped() {
			t.Error("Stopped returned false, but shouldn't")
		}

		s.setStopped(false)

		if s.Stopped() {
			t.Error("Stopped returned true, but shouldn't")
		}
	})

	gracefulServerContext(t, func(s *GracefulServer) {
		time.AfterFunc(20*time.Millisecond, func() {
			s.Stop(0)
		})

		s.ListenAndServe()
		if !s.Stopped() {
			t.Error("Stopped returned false after Stop()")
		}
	})

	gracefulServerContext(t, func(s *GracefulServer) {
		time.AfterFunc(20*time.Millisecond, func() {
			if s.Stopped() {
				t.Error("Server should not be stopped when running")
			}

			s.Stop(0)
		})

		s.ListenAndServe()
	})

	gracefulServerContext(t, func(s *GracefulServer) {
		time.AfterFunc(20*time.Millisecond, func() {
			if s.Stopped() {
				t.Error("Server should not be stopped when running")
			}

			s.Stop(0)
		})

		s.ListenAndServeTLS("foo", "bar")
	})

	gracefulServerContext(t, func(s *GracefulServer) {
		time.AfterFunc(20*time.Millisecond, func() {
			if s.Stopped() {
				t.Error("Server should not be stopped when running")
			}

			s.Stop(0)
		})

		s.ListenAndServeTLSConfig(&tls.Config{})
	})

	gracefulServerContext(t, func(s *GracefulServer) {
		time.AfterFunc(20*time.Millisecond, func() {
			if s.Stopped() {
				t.Error("Server should not be stopped when running")
			}

			s.Stop(0)
		})

		s.Serve(&net.TCPListener{})
	})
}

func ExampleGracefulServer() {
	s := NewGracefulServer(1337,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Foo bar"))
		}))

	// This channel blocks until all connections are closed or the time is up
	sc := s.StopChan()

	// Some go func that stops the server after 2 seconds for no reason
	time.AfterFunc(1*time.Second, func() {
		fmt.Print("Stopping server..")
		s.Stop(10 * time.Second)
	})

	if err := s.ListenAndServe(); err != nil && !s.Stopped() {
		// We didn't stop the server, so something must be wrong
		panic(err)
	}

	// Wait for the server to finish
	<-sc
	fmt.Print("bye!")

	// Output: Stopping server..bye!
}
