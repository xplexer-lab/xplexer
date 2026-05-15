package bus_test

import (
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/xplexer-lab/xplexer/pkg/xkit/bus"
	"sync"
	"testing"
)

func TestRedisBus(t *testing.T) {
	ctx := t.Context()

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
	rBus := bus.NewRedisBus(client)

	t.Cleanup(func() {
		if err := rBus.Close(); err != nil {
			t.Fatalf("failed to close bus: %v", err)
		}
	})

	t.Run("test redis connection", func(t *testing.T) {
		for range 10 {
			err := client.RPush(ctx, "key", "value").Err()
			assert.NoError(t, err)
		}

		size, err := client.LLen(ctx, "key").Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(10), size)
	})

	acceptanceTest(t, rBus)
}

func acceptanceTest(t *testing.T, bus bus.Bus) {
	t.Helper()
	t.Skipf("skip pub sub testing")

	t.Run("pub sub cycle works", func(t *testing.T) {
		t.Skip()

		const msgNr = 5

		var wg sync.WaitGroup
		wg.Add(msgNr * 2)

		go func() {
			for range msgNr {
				//bus.Publish(t.Context())

				wg.Done()
			}
		}()

		//cancel := bus.Subscribe(t.Context(), nil, bus.HandleFn(func() {}))

		wg.Wait()
	})
}
