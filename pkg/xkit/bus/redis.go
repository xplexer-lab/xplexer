package bus

import (
	"context"
	"errors"
	"github.com/cenkalti/backoff/v4"
	"github.com/redis/go-redis/v9"
	"github.com/xplexer-lab/xplexer/pkg/xkit/bus/internal/pb"
	"github.com/xplexer-lab/xplexer/pkg/xkit/errpack"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"log/slog"
	"sync"
	"time"
)

const (
	DefaultPopSize = 100
)

var (
	_ Bus = new(redisBus)
)

type redisBus struct {
	client  *redis.Client
	close   sync.Once
	popSize int
	logger  *slog.Logger
}

func NewRedisBus(client *redis.Client, logger *slog.Logger) Bus {
	return &redisBus{
		client:  client,
		popSize: DefaultPopSize,
		logger:  logger,
	}
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

	err = r.client.
		XAdd(ctx, &redis.XAddArgs{
			Stream: r.stream(envelope.ProtoPayload().ProtoReflect()),
			Values: map[string]interface{}{
				"payload": msg,
				// todo: serialize metadata (timestamp and so on)
			},
		}).
		Err()

	return errpack.Wrap(err, "failed to publish message", errpack.Infra())
}

func (r *redisBus) Subscribe(ctx context.Context, handler AnyHandler) (Cancel, error) {
	ctx, cancel := context.WithCancel(ctx)
	stream := r.stream(handler.getProtoMessage())
	group := handler.Name()

	if err := r.client.XGroupCreateMkStream(ctx, stream, group, "$").Err(); err != nil {
		cancel()
		return nil, err
	}

	go r.subscribeLoop(ctx, handler, stream, group)
	return Cancel(cancel), nil
}

func (r *redisBus) subscribeLoop(
	ctx context.Context,
	handler AnyHandler,
	streamName string,
	groupName string,
) {
	consumerName := "worker-pod-1"

	for {
		id := ">"

		streams, err := r.client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    groupName,
			Streams:  []string{streamName, id},
			Consumer: consumerName,
			Count:    int64(r.popSize),
			NoAck:    false,
			Block:    time.Second * 5,
		}).Result()

		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			continue
		}

		if errors.Is(err, redis.Nil) {
			continue
		}

		if err != nil {
			r.logger.
				With(slog.Any("err", err)).
				ErrorContext(ctx, "failed to pop envelope from redis")
			continue
		}

		pipeline := r.client.Pipeline()

		for _, stream := range streams {
			for _, msg := range stream.Messages {
				if err := handler.Handle(ctx, NewEnvelope(proto.Message(&pb.Person{}))); err != nil {
					r.logger.
						With(slog.Any("err", err)).
						ErrorContext(ctx, "failed to handle message")
					// todo: apply retry strategy
				} else {
					pipeline.XAck(ctx, streamName, groupName, msg.ID)
				}
			}
		}

		if pipeline.Len() > 0 {
			if _, err = pipeline.Exec(ctx); err != nil {
				r.logger.With(slog.Any("err", err)).Error("failed to ack messages in redis")
			}
		}

	}
}

func (r *redisBus) streamByHandler(handler AnyHandler) string {
	return r.stream(handler.getProtoMessage())
}

func (r *redisBus) stream(protoMsg protoreflect.Message) string {
	return string(protoMsg.Descriptor().FullName())
}

func defaultBackoff() backoff.BackOff {
	return backoff.NewExponentialBackOff(
		backoff.WithInitialInterval(time.Millisecond),
		backoff.WithMaxInterval(2*time.Second),
	)
}
