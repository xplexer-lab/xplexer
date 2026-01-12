package api

import (
	"net/http"
	"reflect"
)

type (
	Handler interface {
		http.Handler
		In() reflect.Type
		Out() reflect.Type
	}

	handlerCfg struct {
		errorCodeFallback int
	}

	HandlerOpt func(*handlerCfg)
)

func WithErrorStatusCode(
	code int,
) HandlerOpt {
	return func(hc *handlerCfg) {
		hc.errorCodeFallback = code
	}
}
