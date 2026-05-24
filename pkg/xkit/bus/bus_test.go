package bus_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xplexer-lab/xplexer/pkg/xkit/bus"
	"github.com/xplexer-lab/xplexer/pkg/xkit/bus/internal"
	"sync"
	"testing"
	"time"
)

type busTest struct {
	bus     bus.Bus
	timeout time.Duration
}

func (bt *busTest) run(t *testing.T) {
	t.Helper()
	t.Run("pub sub cycle works", bt.testPubSubCycle)
}

func (bt *busTest) testPubSubCycle(t *testing.T) {
	t.Helper()
	ctx, cancelCtx := context.WithTimeout(t.Context(), bt.timeout)
	t.Cleanup(cancelCtx)

	const msgNr = 50
	var wg sync.WaitGroup
	wg.Add(msgNr * 2)

	cancel, err := bt.bus.Subscribe(ctx, bus.NewAnyHandler[*internal.Person]("handler_name", func(ctx context.Context, e bus.Envelope[*internal.Person]) error {
		defer wg.Done()
		t.Logf("e => %+v", e)
		return nil
	}))

	require.NoError(t, err)

	defer cancel()

	go func() {
		for range msgNr {
			msg := bus.NewEnvelope(&internal.Person{
				Id:       1,
				Email:    "hello@world",
				FullName: "John Doe",
			})

			err := bt.bus.Publish(ctx, msg.AsGeneric())

			assert.NoError(t, err)
			wg.Done()
		}
	}()

	ch := make(chan struct{})
	go func() {
		wg.Wait()
		close(ch)
	}()

	select {
	case <-ch:
		return
	case <-ctx.Done():
		assert.NoError(t, ctx.Err())
	}
}
