package api

import (
	"net/http"
	"reflect"
)

type Handler interface {
	http.Handler
	In() reflect.Type
	Out() reflect.Type
	Path() string
	Method() string
}
