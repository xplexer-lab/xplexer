.PHONY: build

build:
	@go build -o xplexer .

dirs:
	@mkdir -p .tmp

test: dirs
	@go tool gotestsum -f dots-v2 ./...

gen:
	@./scripts/gen.sh
	@mkdir -p ./pkg/xkit/bus/internal/pb
	@protoc  --proto_path=./pkg/xkit/bus --go_out=./pkg/xkit/bus/internal/pb --go_opt=paths=source_relative  dummy.proto
#
wtest:
	@go tool gotestsum --watch --jsonfile=.tmp/test-report.json --post-run-command="go tool xplexer-notify .tmp/test-report.json" ./...

deps:
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	# https://github.com/protocolbuffers/protobuf/releases/download/v35.0/protoc-35.0-linux-x86_64.zip
