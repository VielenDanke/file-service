package middlewares

import "net/http"

// ContentTypeMiddleware ...
type ContentTypeMiddleware struct {
	contentType string
}

// NewContentTypeMiddleware ...
func NewContentTypeMiddleware(contentType string) *ContentTypeMiddleware {
	return &ContentTypeMiddleware{contentType: contentType}
}

// ContentTypeMiddleware ...
func (ctm *ContentTypeMiddleware) ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", ctm.contentType)
		next.ServeHTTP(rw, r)
	})
}
