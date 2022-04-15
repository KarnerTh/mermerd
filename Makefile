GIT_TAG := $(shell git describe --tags `git rev-list --tags --max-count=1`)
test_target := "./..."

.PHONY: test-coverage
test-coverage:
	go test -cover -coverprofile=coverage.out ./...; go tool cover -html=coverage.out -o coverage.html; rm coverage.out

# https://github.com/vektra/mockery is needed
.PHONY: create-mocks
create-mocks:
	mockery --all

# https://github.com/mfridman/tparse is needed
.PHONY: test-all
test-all:
	go test $(test_target) -cover -json | tparse -all

# https://github.com/mfridman/tparse is needed
.PHONY: test-unit
test-unit:
	go test --short $(test_target) -cover -json | tparse -all

.PHONY: publish-package
publish-package:
	echo GOPROXY=proxy.golang.org go list -m github.com/KarnerTh/mermerd@$(GIT_TAG)