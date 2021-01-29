package http

import (
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/unistack-org/micro/v3/api"
)

type handler struct{}

func (h *handler) Method(w http.ResponseWriter, r *http.Request) {

}

func TestRegister(t *testing.T) {
	eps := []*api.Endpoint{
		&api.Endpoint{
			Name:   "service.Method",
			Host:   []string{"example.com"},
			Method: []string{"GET"},
			Path:   []string{"/api/v0/service/{name}"},
		},
		/*
			&api.Endpoint{
				Name:   "service.Call",
				Host:   []string{"example.com"},
				Method: []string{"POST"},
				Path:   []string{"/api/v0/service/{name}"},
			},
		*/
	}

	h := &handler{}
	r := mux.NewRouter()
	if err := Register(r, h, eps); err != nil {
		t.Fatal(err)
	}

	for _, ep := range eps {
		rt := r.Get(ep.Name)
		if rt == nil {
			t.Fatalf("route not registered for %v", ep)
		}
		rx, err := rt.GetPathTemplate()
		if err != nil {
			t.Fatalf("route invalid %v", err)
		}
		if rx != ep.Path[0] {
			t.Fatalf("route invalid %s != %s", rx, ep.Path[0])
		}
	}

}
