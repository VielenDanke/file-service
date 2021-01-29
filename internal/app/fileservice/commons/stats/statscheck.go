// Package stats provides an easy way of publishing of k8s liveness and readiness probes,
// AppVersion and BuildDate, Prometheus metrics.
package stats

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/heptiolabs/healthcheck"
)

type Config struct {
	Version   string
	BuildDate string
	Health    healthcheck.Handler
}

type Option func(*Server)

type Server struct {
	mux *mux.Router
}

func NewServer(options ...Option) *Server {
	res := &Server{
		mux: mux.NewRouter(),
	}

	for _, o := range options {
		o(res)
	}

	return res
}

func NewDefaultServer(conf Config) *Server {
	if conf.Health == nil {
		conf.Health = defaultHealth()
	}

	return NewServer(
		WithHealth(conf.Health),
		WithMetrics(),
		WithVersionDate(conf.Version, conf.BuildDate),
	)
}

func (s Server) Serve(addr string) error {
	return http.ListenAndServe(addr, s.mux)
}

func (s Server) Mux() *mux.Router {
	return s.mux
}

func VersionHandler(version, buildDate string) func(w http.ResponseWriter, r *http.Request) {
	type versionResponse struct {
		Version   string `json:"version"`
		BuildDate string `json:"build_date"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		payload, err := json.Marshal(versionResponse{
			Version:   version,
			BuildDate: buildDate,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("%s", err)))
			return
		}

		_, err = w.Write(payload)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("%s", err)))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
