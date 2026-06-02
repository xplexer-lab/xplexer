package bus

import (
	"fmt"
	"github.com/xplexer-lab/xplexer/pkg/xkit/bus/internal/pb"
	"github.com/xplexer-lab/xplexer/pkg/xkit/errpack"
	"google.golang.org/protobuf/encoding/protojson"
)

type MsgDTO struct {
	Payload  []byte  `json:"payload"`
	Metadata []byte  `json:"metadata"`
	Headers  Headers `json:"headers"`
}

type Headers struct {
	ContentType ContentType `json:"content_type"`
}

func (h *Headers) ToMessage() *pb.Headers {
	return &pb.Headers{
		ContentType: pb.Headers_Unknown,
	}
}

func (h *Headers) getPbContentType() pb.Headers_ContentType {
	switch h.ContentType {
	case ContentTypeUnknown:
		return pb.Headers_Unknown
	case ContentTypeJsonV1:
		return pb.Headers_JSON_V1
	}
	panic(fmt.Errorf("unknown content type: %v", h.ContentType))
}

type ContentType int

const (
	ContentTypeUnknown ContentType = iota
	ContentTypeJsonV1
)

type Serializer interface {
	Serialize(event AnyEnvelope) (MsgDTO, error)
}

type Deserializer interface {
	Deserialize(msg MsgDTO, dst AnyEnvelope) error
}

type serializeDeserializer interface {
	Serializer
	Deserializer
}

var (
	_ serializeDeserializer = (*JsonV1)(nil)
)

type JsonV1 struct {
	headers Headers
}

func NewJsonV1() *JsonV1 {
	return &JsonV1{
		headers: Headers{
			ContentType: ContentTypeJsonV1,
		},
	}
}

func (s *JsonV1) Serialize(event AnyEnvelope) (res MsgDTO, err error) {
	if res.Payload, err = protojson.Marshal(event.ProtoPayload()); err != nil {
		return res, errpack.Wrap(err, "failed to serialize payload", errpack.Domain())
	}

	md := event.Metadata()
	if res.Metadata, err = protojson.Marshal(md.ToProtoMessage()); err != nil {
		return res, errpack.Wrap(err, "failed to serialize payload", errpack.Domain())
	}

	res.Headers = s.headers

	return
}

func (s *JsonV1) Deserialize(msg MsgDTO, dest AnyEnvelope) error {
	if msg.Headers.ContentType != ContentTypeJsonV1 {
		return errpack.New("invalid content type passed to JsonV1 deserializer", errpack.Domain())
	}

	if err := protojson.Unmarshal(msg.Payload, dest.ProtoPayload()); err != nil {
		return errpack.New("failed to deserialize message payload", errpack.Infra())
	}

	md := &pb.Metadata{}
	if err := protojson.Unmarshal(msg.Metadata, md); err != nil {
		return errpack.New("failed to deserialize metadata", errpack.Infra())
	}

	if err := dest.Metadata().FromProtoMessage(md); err != nil {
		return errpack.New("failed to deserialize metadata", errpack.Infra())
	}

	return nil
}
