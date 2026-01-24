package usecases

import (
	"context"

	"github.com/xplexer-lab/xplexer/internal/common/restapi"
)

type ProjectListIn struct{}

type ProjectDto struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

var projectList = restapi.Operation(func(ctx context.Context, in ProjectListIn) ([]ProjectDto, error) {
	return []ProjectDto{
		{"1", "first"},
		{"2", "second"},
	}, nil
})
