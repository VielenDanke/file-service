package stats

import (
	"net/http/pprof"

	"github.com/heptiolabs/healthcheck"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func WithHealth(h healthcheck.Handler) Option {
	return func(s *Server) {
		s.mux.HandleFunc("/live", h.LiveEndpoint)
		s.mux.HandleFunc("/ready", h.ReadyEndpoint)
	}
}

func WithDefaultHealth() Option {
	return WithHealth(defaultHealth())
}

func WithVersionDate(version, date string) Option {
	return func(s *Server) {
		s.mux.HandleFunc("/version", VersionHandler(version, date))
	}
}

func WithMetrics() Option {
	return func(s *Server) {
		s.mux.Handle("/metrics", promhttp.Handler())
	}
}

func WithProfile() Option {
	return func(s *Server) {
		s.mux.HandleFunc("/debug/pprof/", pprof.Index)
		s.mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		s.mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		s.mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		s.mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
		s.mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
		s.mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
		s.mux.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
		s.mux.Handle("/debug/pprof/block", pprof.Handler("block"))
		s.mux.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
		s.mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	}
}

func defaultHealth() healthcheck.Handler {
	health := healthcheck.NewHandler()

	health.AddLivenessCheck("basic live", func() error {
		return nil
	})
	health.AddReadinessCheck("basic ready", func() error {
		return nil
	})

	return health
}
