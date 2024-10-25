lint: build
	golangci-lint run ./...

test: build
	go clean --testcache
	go run gotest.tools/gotestsum@latest --format testname

build:
	go build ./...