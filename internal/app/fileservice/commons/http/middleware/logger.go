package middleware

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/metadata"
)

type LoggerMiddleware struct{}

func NewLoggerMiddleware() *LoggerMiddleware {
	return &LoggerMiddleware{}
}

func (s *LoggerMiddleware) Wrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		switch r.URL.Path {
		case "/live", "/ready", "/metrics", "/version":
			next.ServeHTTP(w, r)
			return
		}

		log, ok := logger.FromContext(r.Context())
		if !ok {
			log = logger.DefaultLogger
		}
		var body []byte
		if r.Body != nil {
			// use http.MaxBytesReader to avoid DoS
			body, _ = ioutil.ReadAll(r.Body)
			// Restore the io.ReadCloser to its original state
			r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		}
		rw := httptest.NewRecorder()
		next.ServeHTTP(rw, r)

		log.Fields(map[string]interface{}{
			"http_method": r.Method,
			"http_uri":    r.URL.String(),
			//		"http_reqbody": strconv.Quote(string(body)),
			//		"http_rspbody": strconv.Quote(string(rw.Body.Bytes())),
			"http_code": rw.Code,
		}).Info(context.Background())
		// this copies the recorded response to the response writer
		for k, v := range rw.HeaderMap {
			w.Header()[k] = v
		}
		if _, ok := rw.HeaderMap[MetadataKey]; !ok {
			if id, ok := metadata.Get(r.Context(), MetadataKey); ok {
				w.Header()[MetadataKey] = []string{id}
			}
		}

		w.WriteHeader(rw.Code)
		rw.Body.WriteTo(w)
	})
}
