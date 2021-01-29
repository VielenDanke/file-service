package http

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
	"github.com/unistack-org/micro/v3/api"
)

var (
	ErrInvalidHandler = errors.New("invalid handler type")
)

func Register(r *mux.Router, h interface{}, eps []*api.Endpoint) error {

	v := reflect.ValueOf(h)

	methods := v.NumMethod()
	if methods < 1 {
		return ErrInvalidHandler
	}

	for _, ep := range eps {
		idx := strings.Index(ep.Name, ".")
		if idx < 1 || len(ep.Name) <= idx {
			return fmt.Errorf("invalid api.Endpoint name: %s", ep.Name)
		}
		name := ep.Name[idx+1:]
		m := v.MethodByName(name)
		if !m.IsValid() || m.IsZero() {
			return fmt.Errorf("invalid handler, method %s not found", name)
		}

		rh, ok := m.Interface().(func(http.ResponseWriter, *http.Request))
		if !ok {
			return fmt.Errorf("invalid handler: %#+v", m.Interface())
		}
		r.HandleFunc(ep.Path[0], rh).Methods(ep.Method...).Name(ep.Name)
	}

	return nil
}
