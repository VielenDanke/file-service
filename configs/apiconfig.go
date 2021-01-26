package configs

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
	"github.com/unistack-org/micro/v3/api"
)

// ConfigureHandlerToEndpoints ...
func ConfigureHandlerToEndpoints(router *mux.Router, handler interface{}, endpoints []*api.Endpoint) error {
	for _, v := range endpoints {
		fName := strings.Split(v.Name, ".")[1]
		f, ok := reflect.ValueOf(handler).MethodByName(fName).Interface().(func(w http.ResponseWriter, r *http.Request))
		if !ok {
			return fmt.Errorf("Naming of method is incorrect: %v:%v", fName, v.Path)
		}
		for k, j := range v.Path {
			router.HandleFunc(j, f).Methods(v.Method[k])
		}
	}
	return nil
}
