package cli

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/xplexer-lab/xplexer/internal/common/errpack"
	"github.com/xplexer-lab/xplexer/internal/usecases"
)

type logLevel string

const (
	logLevelDebug logLevel = "debug"
	logLevelInfo  logLevel = "info"
	logLevelWarn  logLevel = "warn"
	logLevelError logLevel = "error"
)

const (
	helpRoot = `
	TODO: add some base help
	`
)

func New(name string) *Cli {

	cmd := Cli{
		root: kingpin.New(name, helpRoot),
	}
	cmd.logger = cmd.buildLogger()

	cmd.root.Flag("log-level", "Set logging level").Default("debug").EnumVar(&cmd.logLevel, "debug", "info", "warn", "error")
	cmd.root.PreAction(func(_ *kingpin.ParseContext) error {
		// Rebuild logger because of severity level depends on flag
		cmd.logger = cmd.buildLogger()
		return nil
	})

	cmd.workerCommand()
	cmd.controlPlaneCommand()
	cmd.controlPlaneUiCommand()
	cmd.configCommand()

	return &cmd
}

type Cli struct {
	root     *kingpin.Application
	logger   *slog.Logger
	logLevel string
}

func (c *Cli) Run(args []string) {
	if command, err := c.root.Parse(args); err != nil {
		c.logger.Error("failed to execute", slog.Any("err", err))
	} else {
		c.logger.Debug("executed", slog.String("command", command))
	}
}

func (c *Cli) workerCommand() {
	worker := c.root.Command("worker", "run worker")

	worker.Action(func(pc *kingpin.ParseContext) error {
		return errpack.New("not implemented", errpack.Bootstrap())
	})
}

func (c *Cli) getLoggerLevel() slog.Level {
	switch logLevel(c.logLevel) {
	case logLevelError:
		return slog.LevelError
	case logLevelWarn:
		return slog.LevelWarn
	case logLevelInfo:
		return slog.LevelInfo
	case logLevelDebug:
		fallthrough
	default:
		return slog.LevelDebug
	}
}

func (c *Cli) buildLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: c.getLoggerLevel(),
	}))
}

func (c *Cli) controlPlaneCommand() {
	controlPlane := c.root.Command("control-plane", "contol plane")

	config := controlPlane.Flag("config", "Path to configuration file").Short('c').String()
	port := controlPlane.Flag("port", "Listen port").Short('p').Default("8080").Uint16()
	host := controlPlane.Flag("host", "Listen host").Short('h').Default("").String()

	controlPlane.Action(func(pc *kingpin.ParseContext) error {
		c.logger.Info("reading config file", slog.String("path", *config))

		handler, err := usecases.BuildRouter().BuildHandler()

		if err != nil {
			return err
		}

		listen := fmt.Sprintf("%s:%d", *host, *port)

		c.logger.Debug("starting server", slog.String("listen", listen))

		// todo: todo: graceful handler
		http.ListenAndServe(
			listen,
			handler,
		)

		return nil
	})
}

func (c *Cli) controlPlaneUiCommand() {
	cp := c.root.Command("control-plane-ui", "contol plane ui interface")
	cp.Action(func(pc *kingpin.ParseContext) error {
		return errpack.New("non implemented", errpack.Bootstrap())
	})
}

func (c *Cli) configCommand() {
	cp := c.root.Command("config", "config related actions")

	cp.
		Command("validate", "validate config").
		Arg("path", "path to condif file to validate").
		Required().
		Action(func(pc *kingpin.ParseContext) error {
			return errpack.New("non implemented", errpack.Bootstrap())
		})

	cp.
		Command("gen", "generate example file name").
		Action(func(pc *kingpin.ParseContext) error {
			return errpack.New("not implemented", errpack.Bootstrap())
		})
}
