GIT_TAG := $(shell git describe --tags --abbrev=0)
test_target := "./..."

.PHONY: prepare-sqlite
prepare-sqlite:
	rm -f mermerd_test.db
	sqlite3 mermerd_test.db < <(cat test/db-table-setup.sql test/sqlite/sqlite-enum-setup.sql test/sqlite/sqlite-multiple-databases.sql)

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

.PHONY: test-cleanup
test-cleanup:
	go clean -testcache

.PHONY: publish-package
publish-package:
	GOPROXY=proxy.golang.org go list -m github.com/KarnerTh/mermerd@$(GIT_TAG)
