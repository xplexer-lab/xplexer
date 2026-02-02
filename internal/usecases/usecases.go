package usecases

import (
	"log/slog"
	"os"

	"github.com/xplexer-lab/xplexer/pkg/xkit/consistency"
	"github.com/xplexer-lab/xplexer/pkg/xkit/logger"
	"github.com/xplexer-lab/xplexer/pkg/xkit/restapi"
)

func BuildRouter(txm TxManager) *restapi.Router {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	router := restapi.NewRouter()
	router.Use(logger.Inject(log))
	router.Use(consistency.Inject(txm))
	router.SetLogger(log)

	router.Post("/core/projects", projectCreate)
	router.Get("/core/projects", projectList)
	router.Get("/core/projects/{id}", projectGet)

	return router
}
