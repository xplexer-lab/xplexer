package domain

import "github.com/xplexer-lab/xplexer/pkg/xkit/entity"

type Tenant struct {
	*entity.Entity
	name string
}
