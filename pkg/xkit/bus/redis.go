package bus

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/xplexer-lab/xplexer/pkg/xkit/errpack"
	"google.golang.org/protobuf/encoding/protojson"
	"sync"
)

var (
	_ Bus = new(redisBus)
)

type redisBus struct {
	client *redis.Client
	close  sync.Once
}

func NewRedisBus(client *redis.Client) Bus {
	return &redisBus{client: client}
}

func (r *redisBus) Close() error {
	var err error

	r.close.Do(func() {
		err = errpack.Wrap(
			r.client.Close(),
			"failed to close redis",
			errpack.Infra(),
		)
	})

	return err
}

func (r *redisBus) Publish(ctx context.Context, envelope AnyEnvelope) error {
	msg, err := protojson.Marshal(envelope.ProtoPayload())

	if err != nil {
		return errpack.Wrap(err, "failed to encode envelope")
	}

	return r.client.RPush(ctx, "message", msg).Err()
}

func (r *redisBus) Subscribe(ctx context.Context, envelope AnyEnvelope, handler AnyHandler) Cancel {
	return func() {}
}
