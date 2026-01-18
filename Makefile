.PHONY: build

build: deps
	go build -o xplexer .

dirs:
	mkdir -p .tmp

test: deps dirs
	gotestsum -f dots ./...

deps:
	go install gotest.tools/gotestsum@v1.13

