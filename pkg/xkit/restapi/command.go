package restapi

import (
	"context"
	"net/http"
)

type (
	CommandOpt func(*command)

	command struct {
		operationHanlder
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
	opts ...OperationOpt,
) Handler {

	defaultOpts := []OperationOpt{
		WithSuccessCode(http.StatusAccepted),
	}

	return Operation(
		append(defaultOpts, opts...)...,
	)
}

// CommandHandler
// Defines handler for the command.
func CommandHandler[In any](handle func(context.Context, In) error) OperationOpt {
	return WithHandler(func(ctx context.Context, in In) (commandOut, error) {
		if err := handle(ctx, in); err != nil {
			return commandOut{}, err
		}

		return commandOut{
			Message: defaultCommandMessage,
		}, nil
	})
}
