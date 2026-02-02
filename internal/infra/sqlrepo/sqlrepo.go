package sqlrepo

import (
	"context"
	"database/sql"

	"github.com/xplexer-lab/xplexer/internal/usecases"
	"github.com/xplexer-lab/xplexer/pkg/xkit/consistency"
	"github.com/xplexer-lab/xplexer/pkg/xkit/errpack"
	"gorm.io/gorm"
)

func NewTxManager(dialector gorm.Dialector, cfg *gorm.Config) (usecases.TxManager, error) {
	db, err := gorm.Open(dialector, cfg)

	if err != nil {
		return nil, errpack.Wrap(err, "failed to create gorm connection", errpack.Bootstrap())
	}

	if err := db.AutoMigrate(&projectModel{}); err != nil {
		return nil, errpack.Wrap(err, "failed to migrate", errpack.Bootstrap())
	}

	return &txManager{
		db: db,
	}, nil
}

var (
	_ usecases.TxManager = new(txManager)
)

type txManager struct {
	db *gorm.DB
}

func (tx *txManager) Do(
	ctx context.Context,
	handle func(context.Context, usecases.Repositories) (any, error),
	_ ...consistency.TxOpt,
) (any, error) {
	var res any

	err := tx.db.
		WithContext(ctx).
		Transaction(
			func(tx *gorm.DB) error {
				var err error
				res, err = handle(ctx, newReposiotries(tx))
				return err
			},
			&sql.TxOptions{Isolation: sql.LevelDefault}, // todo: add possibility to change isolation level
		)

	return res, err
}

func newReposiotries(tx *gorm.DB) usecases.Repositories {
	return usecases.Repositories{
		Projects: newProjectRepository(tx),
	}
}
