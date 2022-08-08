.PHONY: lint
lint:
	@hash golangci-lint > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.39.0; \
	fi
	golangci-lint run

.PHONY: test
test:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./bittrex

.PHONY: test-examples
test-examples:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./examples

.PHONY: test-examples-http
test-examples-http:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./examples -run TestFlagsHttp

.PHONY: test-examples-ws
test-examples-ws:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./examples -run TestFlagsWs

.PHONY: test-all
test-all:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...