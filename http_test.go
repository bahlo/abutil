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

func TestRemoteIp(t *testing.T) {
	mockRequestContext(t, func(r *http.Request) {
		ip := "123.456.7.8"

		r.Header = http.Header{
			"X-Real-Ip": []string{ip},
		}

		out := RemoteIp(r)
		if out != ip {
			t.Errorf("Expected %s, but got %s", ip, out)
		}
	})
}

func TestRemoteIpForwardedFor(t *testing.T) {
	mockRequestContext(t, func(r *http.Request) {
		ip := "123.456.7.8"

		r.Header = http.Header{
			"X-Forwarded-For": []string{ip},
		}

		out := RemoteIp(r)
		if out != ip {
			t.Errorf("Expected %s, but got %s", ip, out)
		}
	})
}

func TestRemoteIpRemoteAddr(t *testing.T) {
	mockRequestContext(t, func(r *http.Request) {
		ip := "123.456.7.8"
		r.RemoteAddr = ip

		out := RemoteIp(r)
		if out != ip {
			t.Errorf("Expected %s, but got %s", ip, out)
		}
	})
}

func TestRemoteIpLocalhost(t *testing.T) {
	mockRequestContext(t, func(r *http.Request) {
		ip := "127.0.0.1"
		r.RemoteAddr = "["

		out := RemoteIp(r)
		if out != ip {
			t.Errorf("Expected %s, but got %s", ip, out)
		}
	})
}

func remoteIpMockServe(h http.HandlerFunc) {
	mockRequestContext(nil, func(r *http.Request) {
		r.RemoteAddr = "123.456.7.8"

		w := httptest.NewRecorder()
		h(w, r)
	})
}

func ExampleRemoteIp() {
	someHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("New request from %s\n", RemoteIp(r))
	}

	// Get's called with RemoteAddr = 123.456.7.8
	remoteIpMockServe(someHandler)

	// Output: New request from 123.456.7.8
}
