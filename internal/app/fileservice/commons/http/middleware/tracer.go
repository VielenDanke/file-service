package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/opentracing-contrib/go-gorilla/gorilla"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/unistack-org/micro/v3/metadata"
)

func NewTracerMiddleware(router *mux.Router, tracer opentracing.Tracer) mux.MiddlewareFunc {
	opNameFunc := func(r *http.Request) string {
		if route := mux.CurrentRoute(r); route != nil {
			return route.GetName()
		} else {
			var m mux.RouteMatch
			if router.Match(r, &m) && m.Route != nil {
				return m.Route.GetName()
			}
		}
		return "unknown"
	}

	spanFilterFunc := func(r *http.Request) bool {
		switch r.URL.Path {
		case "/live", "/ready", "/metrics", "/version":
			return false
		}
		switch r.Method {
		case http.MethodHead, http.MethodOptions:
			return false
		}
		return true
	}

	observerFunc := func(span opentracing.Span, r *http.Request) {
		id, _ := metadata.Get(r.Context(), MetadataKey)
		span.SetTag(MetadataKey, id)
	}

	tracefunc := func(next http.Handler) http.Handler {
		ret := gorilla.Middleware(tracer, next,
			[]nethttp.MWOption{
				nethttp.MWComponentName("apigw"),
				nethttp.MWSpanFilter(spanFilterFunc),
				nethttp.OperationNameFunc(opNameFunc),
				nethttp.MWSpanObserver(observerFunc),
			}...,
		)
		return ret
	}

	return tracefunc

}
