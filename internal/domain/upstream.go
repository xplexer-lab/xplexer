package domain

import "github.com/xplexer-lab/xplexer/pkg/xkit/entity"

type Upstream struct {
	*entity.Entity
	project entity.Id
}
