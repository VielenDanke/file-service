package handlers

import (
	"net/http"

	"github.com/unistack-org/micro/v3/codec"
)

// SomeHandler ...
type SomeHandler struct {
	codec codec.Codec
}

// NewSomeHandler ...
func NewSomeHandler(codec codec.Codec) *SomeHandler {
	return &SomeHandler{
		codec: codec,
	}
}

// SomeSave ...
func (sh *SomeHandler) SomeSave(w http.ResponseWriter, r *http.Request) {

}
