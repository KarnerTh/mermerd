GIT_TAG := $(shell git describe --tags --abbrev=0)
test_target := "./..."

.PHONY: prepare-sqlite
prepare-sqlite:
	rm -f mermerd_test.db
	cat test/db-table-setup.sql test/sqlite/sqlite-enum-setup.sql test/sqlite/sqlite-multiple-databases.sql | sqlite3 mermerd_test.db

.PHONY: test-coverage
test-coverage:
	go test -cover -coverprofile=coverage.out ./...; go tool cover -html=coverage.out -o coverage.html; rm coverage.out

# https://github.com/vektra/mockery is needed
.PHONY: create-mocks
create-mocks:
	mockery --all

.PHONY: test-all
test-all: test-setup
	go test $(test_target) -cover -json | tparse -all

.PHONY: test-unit
test-unit: test-setup
	go test --short $(test_target) -cover -json | tparse -all

.PHONY: test-setup
test-setup:
	go install github.com/mfridman/tparse@v0.18.0
	cd test && docker-compose up -d

.PHONY: test-cleanup
test-cleanup:
	go clean -testcache
	cd test && docker-compose stop && docker-compose rm -f
	rm mermerd_test.db 2&> /dev/null || true

.PHONY: publish-package
publish-package:
	GOPROXY=proxy.golang.org go list -m github.com/KarnerTh/mermerd@$(GIT_TAG)
