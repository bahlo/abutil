package abutil

import (
	"net/http"
	"strings"
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
