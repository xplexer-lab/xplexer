package restapi

import (
	"context"
	"net/http"
	"slices"
)

type (
	CommandOpt[In any] func(*command[In])

	command[In any] struct {
		queryHandler[In, commandOut]
		message string
	}

	commandOut struct {
		Message string `json:"message"`
	}
)

const (
	defaultCommandMessage = "Accepted"
)

// Command
// Creates command handler for the router.
// Command is a command in a context of `CQRS`.
// Command can have defined input `DTO` and response code.
// Generally command handler responds with status code 202 Accepted.
func Command[In any](
	handle func(context.Context, In) error,
	opts ...QueryOpt[In, commandOut],
) Handler {
	defaultOpts := []QueryOpt[In, commandOut]{
		WithQueryCommon[In, commandOut](
			WithSuccessCode(http.StatusAccepted),
		),
	}

	return Query(func(ctx context.Context, in In) (commandOut, error) {

		if err := handle(ctx, in); err != nil {
			return commandOut{}, err
		}

		return commandOut{
			Message: defaultCommandMessage,
		}, nil
	}, slices.Concat(defaultOpts, opts)...)
}
