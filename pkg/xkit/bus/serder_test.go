package bus_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/xplexer-lab/xplexer/pkg/xkit/bus"
	"github.com/xplexer-lab/xplexer/pkg/xkit/bus/internal/pb"
	"testing"
)

func Test_JsonV1(t *testing.T) {

	serializer := bus.NewJsonV1()

	// todo: generate acceptance tests for each type of serializer
	t.Run("serialize/deserialize cycle works", func(t *testing.T) {
		msg := bus.NewEnvelope(&pb.Person{
			Id:       1234,
			FullName: "John Doe",
			Email:    "john.doe@example.com",
		})

		dto, err := serializer.Serialize(&msg)
		assert.NoError(t, err)

		deserialized := bus.NewEnvelope(&pb.Person{})
		err = serializer.Deserialize(dto, &deserialized)
		assert.NoError(t, err)

		deserializedPayload := deserialized.Payload()

		assert.Equal(t, int32(1234), deserializedPayload.Id)
		assert.Equal(t, "John Doe", deserializedPayload.FullName)
		assert.Equal(t, "john.doe@example.com", deserializedPayload.Email)
		assert.Equal(t, msg.Payload(), deserializedPayload)
		assert.Equal(t, msg.Metadata(), deserialized.Metadata())
	})
}
