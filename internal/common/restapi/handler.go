package restapi

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
		successCode       int
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

func WithSuccessCode(code int) HandlerOpt {
	return func(hc *handlerCfg) {
		hc.successCode = code
	}
}

func (hc handlerCfg) getSuccessCode() int {
	if hc.successCode > 0 {
		return hc.successCode
	} else {
		return http.StatusOK
	}
}

func (hc handlerCfg) getErrorCodeFallback() int {
	if hc.errorCodeFallback > 0 {
		return hc.errorCodeFallback
	} else {
		return http.StatusInternalServerError
	}
}
