.PHONY: build

build:
	@go build -o xplexer .

dirs:
	@mkdir -p .tmp

test: dirs
	@go tool gotestsum -f dots-v2 ./...

gen:
	./scripts/gen.sh
#
wtest:
	@go tool gotestsum --watch --jsonfile=.tmp/test-report.json --post-run-command="go tool xplexer-notify .tmp/test-report.json" ./...
