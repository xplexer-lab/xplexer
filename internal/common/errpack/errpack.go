package errpack

import "errors"

var (
	_ error = new(Error)
)

type (
	Error struct {
		msg  string
		prev error
		typ  Type
	}

	Type struct {
		val string
	}

	Opt func(*Error)
)

var (
	TypeUnknown      = Type{"unknown"}
	TypeDomain       = Type{"domain"}
	TypeInfra        = Type{"infra"}
	TypeBootstrap    = Type{"bootstrap"}
	TypeUnauthorized = Type{"unauthroized"}
	TypeForbidden    = Type{"forbidden"}
	TypeUnreachable  = Type{"unreachable"}
)

func New(msg string, opts ...Opt) *Error {
	err := &Error{
		msg: msg,
		typ: TypeUnknown,
	}

	for _, apply := range opts {
		if apply != nil {
			apply(err)
		}
	}

	return err
}

func Wrap(err error, msg string, opts ...Opt) error {
	if err == nil {
		return nil
	}

	return New(msg, append(opts, WithPrev(err))...)
}

func (e *Error) Error() string {
	return e.msg
}

func (e *Error) Type() Type {
	return e.typ
}

func (e Error) Wrap(prev error) *Error {
	e.prev = prev
	return &e
}

func (e *Error) Unwrap() error {
	return e.prev
}

func (e *Error) Is(err error) bool {
	var perr = new(Error)
	return errors.As(err, &perr) && perr.typ == e.typ && perr.msg == e.msg
}

func withType(typ Type) Opt {
	return func(err *Error) {
		err.typ = typ
	}
}

func Domain() Opt {
	return withType(TypeDomain)
}

func Infra() Opt {
	return withType(TypeInfra)
}

func Bootstrap() Opt {
	return withType(TypeBootstrap)
}

func Forbidden() Opt {
	return withType(TypeForbidden)
}

func Unauthorized() Opt {
	return withType(TypeUnauthorized)
}

func Unreachable() Opt {
	return withType(TypeUnreachable)
}

func WithPrev(prev error) Opt {
	return func(err *Error) {
		err.prev = prev
	}
}
