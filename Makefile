.PHONY: build

build: deps
	@go build -o xplexer .

dirs:
	@mkdir -p .tmp

test: deps dirs
	@gotestsum -f dots-v2 ./...

# tools
deps:
	@go install gotest.tools/gotestsum@v1.13
	@go install ./cmd/xplexer-gen-query
	@go generate

