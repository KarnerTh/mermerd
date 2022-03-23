test_target := "./..."

.PHONY: test-coverage
test-coverage:
	go test -cover -coverprofile=coverage.out ./...; go tool cover -html=coverage.out -o coverage.html; rm coverage.out

# https://github.com/mfridman/tparse
.PHONY: test-pretty
test-pretty:
	go test $(test_target) -cover -json | tparse -all

.PHONY: test-pretty-small
test-pretty-small:
	go test $(test_target) -cover -json | tparse -all -smallscreen
