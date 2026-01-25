.PHONY: build

build: deps
	@go build -o xplexer .

dirs:
	@mkdir -p .tmp

test: deps dirs
	@go tool gotestsum -f dots-v2 ./...

wtest:
	@go tool gotestsum --watch --jsonfile=.tmp/test-report.json --post-run-command="go tool xplexer-notify .tmp/test-report.json" ./...
