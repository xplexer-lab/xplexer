package domain

import "github.com/xplexer-lab/xplexer/pkg/xkit/entity"

type Message struct {
	*entity.Entity
	project       entity.Id
	schema        SchemaSlug
	schemaVersion uint
	payload       []byte
}
