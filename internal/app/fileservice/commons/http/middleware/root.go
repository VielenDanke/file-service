package middleware

import (
	"net/http"
)

type RootMiddleware struct{}

func (s *RootMiddleware) Wrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if r.RequestURI == "/" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
