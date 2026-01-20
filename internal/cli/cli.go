package cli

import (
	"log/slog"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/xplexer-lab/xplexer/internal/common/errpack"
)

const (
	helpRoot = `
	TODO: add some base help
	`
)

func New(name string) *Cli {

	cmd := Cli{
		root:    kingpin.New(name, helpRoot),
		verbose: false,
	}
	cmd.logger = cmd.buildLogger()

	cmd.root.Flag("verbose", "Enable verbose mode").Short('v').BoolVar(&cmd.verbose)
	cmd.root.PreAction(func(_ *kingpin.ParseContext) error {
		// Rebuild logger because of severity level depends on flag
		cmd.logger = cmd.buildLogger()
		return nil
	})

	cmd.workerCommand()
	cmd.cpCommand()
	cmd.configCommand()

	return &cmd
}

type Cli struct {
	root    *kingpin.Application
	logger  *slog.Logger
	verbose bool
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

	config := worker.Flag("config", "Path to configuration file").Short('c').String()

	worker.Action(func(pc *kingpin.ParseContext) error {
		c.logger.Info("reading config file", slog.String("path", *config))
		return nil
	})
}

func (c *Cli) buildLogger() *slog.Logger {
	level := slog.LevelWarn

	if c.verbose {
		level = slog.LevelDebug
	}

	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
}

func (c *Cli) cpCommand() {
	cp := c.root.Command("cp", "contol plane")
	cp.Action(func(pc *kingpin.ParseContext) error {
		return errpack.New("non implemented", errpack.WithBootstrap())
	})
}

func (c *Cli) configCommand() {
	cp := c.root.Command("config", "config related actions")

	cp.
		Command("validate", "validate config").
		Action(func(pc *kingpin.ParseContext) error {
			return errpack.New("non implemented", errpack.WithBootstrap())
		})

	cp.
		Command("gen", "generate example file name").
		Action(func(pc *kingpin.ParseContext) error {
			return errpack.New("not implemented", errpack.WithBootstrap())
		})
}
