package api

var (
	_ Handler = new(queryHandler)
)

type (
	queryHandler struct {
	}
)

func Query[In, Out any]() Handler {
	return nil
}
