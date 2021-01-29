package middleware

import (
	"net/http"
	"net/textproto"

	"github.com/google/uuid"
	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/metadata"
)

var (
	MetadataKey = textproto.CanonicalMIMEHeaderKey("x-request-id")
	loggerKey   = "x-request-id"
)

type RequestIDMiddleware struct{}

func NewRequestIDMiddleware() *RequestIDMiddleware {
	return &RequestIDMiddleware{}
}

func (s *RequestIDMiddleware) Wrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(MetadataKey)
		if id == "" {
			uid, err := uuid.NewRandom()
			if err != nil {
				uid = uuid.Nil
			}
			id = uid.String()
		}
		ctx := metadata.Set(r.Context(), MetadataKey, id)
		ctx = logger.NewContext(ctx, logger.DefaultLogger.Fields(map[string]interface{}{loggerKey: id}))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
