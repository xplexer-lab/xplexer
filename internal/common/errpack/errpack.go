package errpack

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
		apply(err)
	}

	return err
}

func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}

	return &Error{
		msg:  msg,
		prev: err,
	}
}

func (e *Error) Error() string {
	return e.msg
}

func (e *Error) Type() Type {
	return e.typ
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
