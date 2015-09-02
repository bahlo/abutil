package abutil

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
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
