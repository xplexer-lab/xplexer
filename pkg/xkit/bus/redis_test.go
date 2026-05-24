package bus_test

import (
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/xplexer-lab/xplexer/pkg/xkit/bus"
	"log/slog"
	"testing"
	"time"
)

func TestRedisBus(t *testing.T) {
	ctx := t.Context()
	discardLogger := slog.New(slog.DiscardHandler)

	redisC, err := testcontainers.Run(
		ctx, "redis:latest",
		testcontainers.WithExposedPorts("6379/tcp"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("6379/tcp"),
			wait.ForLog("Ready to accept connections"),
		),
	)
	testcontainers.CleanupContainer(t, redisC)
	require.NoError(t, err)

	endpoint, err := redisC.Endpoint(ctx, "")
	require.NoError(t, err)
	t.Logf("endpoint=%s", endpoint)

	client := redis.NewClient(&redis.Options{
		Addr: endpoint,
	})

	rBus := bus.NewRedisBus(client, discardLogger)

	t.Cleanup(func() {
		if err := rBus.Close(); err != nil {
			t.Fatalf("failed to close bus: %v", err)
		}
	})

	t.Run("test redis client acceptance tests", func(t *testing.T) {
		at := &busTest{bus: rBus, timeout: time.Second * 300}
		at.run(t)
	})
}
