package bus

import "context"

type Publisher interface {
	Publish(context.Context, AnyEnvelope)
}
