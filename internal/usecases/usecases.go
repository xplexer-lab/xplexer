package usecases

import (
	"log/slog"
	"os"

	"github.com/xplexer-lab/xplexer/internal/common/logger"
	"github.com/xplexer-lab/xplexer/internal/common/restapi"
)

func BuildRouter() *restapi.Router {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	router := restapi.NewRouter()
	router.Use(logger.Inject(log))
	router.SetLogger(log)

	router.Post("/core/projects", projectCreate)
	router.Get("/core/projects", projectList)

	return router
}
