package domain

import "github.com/xplexer-lab/xplexer/internal/common/entity"

type Tenant struct {
	*entity.Entity
	name string
}
