package middleware

import (
	"net/http"
	"time"
)

var (
	// Unix epoch time

	epoch = time.Unix(0, 0).Format(time.RFC1123)

	// Taken from https://github.com/mytrile/nocache
	noCacheHeaders = map[string]string{
		"Expires":         epoch,
		"Cache-Control":   "no-cache, no-store, no-transform, must-revalidate, private, max-age=0",
		"Pragma":          "no-cache",
		"X-Accel-Expires": "0",
	}

	etagHeaders = []string{
		"ETag",
		"If-Modified-Since",
		"If-Match",
		"If-None-Match",
		"If-Range",
		"If-Unmodified-Since",
	}
)

type NocacheMiddleware struct{}

func NewNocacheMiddleware() *NocacheMiddleware {
	return &NocacheMiddleware{}
}

func nocacheHaders(w http.ResponseWriter, r *http.Request) {
	// Delete any ETag headers that may have been set
	for _, v := range etagHeaders {
		if r.Header.Get(v) != "" {
			r.Header.Del(v)
		}
	}

	// Set our NoCache headers
	for k, v := range noCacheHeaders {
		w.Header().Set(k, v)
	}
}

func (s *NocacheMiddleware) Wrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nocacheHaders(w, r)
		next.ServeHTTP(w, r)
	})
}
