package restapi

type (
	CommandOpt[In any] func(*command[In])

	command[T any] struct {
	}
)

// Command
// Creates command handler for the router.
// Command is a command in a context of `CQRS`.
// Command can have defined input `DTO` and response code.
// Generally command handler responds with status code 202 Accepted.
func Command[In any](
	handle func(Context, In) error,
	opts ...CommandOpt[In],
) Handler {
	panic("not implemented")
}
