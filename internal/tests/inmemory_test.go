package tests_test

import (
	"testing"

	"github.com/xplexer-lab/xplexer/internal/tests"
)

func TestInmemory(t *testing.T) {
	tests.New().Run(t)
}

