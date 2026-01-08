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
	Unknown = Type{"unknown"}
	Domain  = Type{"domain"}
	Infra   = Type{"infra"}
)

func New(msg string, opts ...Opt) *Error {
	err := &Error{
		msg: msg,
		typ: Unknown,
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

func WithDomain() Opt {
	return withType(Domain)
}

func WithInfra() Opt {
	return withType(Infra)
}

func WithPrev(prev error) Opt {
	return func(err *Error) {
		err.prev = prev
	}
}
