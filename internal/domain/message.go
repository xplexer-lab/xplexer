package domain

import "github.com/xplexer-lab/xplexer/internal/common/entity"

type Message struct {
	*entity.Entity
	project entity.Id
}
