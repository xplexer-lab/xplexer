package main

import (
	"os"

	"github.com/xplexer-lab/xplexer/internal/cli"
)

func main() {
	cli.New(os.Args[0]).Run(os.Args[1:])
}
