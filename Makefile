.PHONY: build

JUNIT_TEST_FILE=.tmp/test-ci.junit.xml

build: deps
	go build -o xplexer .

dirs:
	mkdir -p .tmp

test: deps dirs
	gotestsum -f dots ./...

test-ci: deps dirs
	gotestsum -f github-actions --junitfile=${JUNIT_TEST_FILE} ./...
	@echo "JUNIT_TEST_FILE"=${JUNIT_TEST_FILE}

deps:
	go install gotest.tools/gotestsum@v1.13

tmpdir:
	@echo -n ${JUNIT_TEST_FILE}
