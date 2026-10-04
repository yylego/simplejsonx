COVERAGE_DIR ?= .coverage.out

.PHONY: generate test test-with-flags

generate:
	@go generate ./...

test:
	@mkdir -p "$(COVERAGE_DIR)"
	make test-with-flags TEST_FLAGS='-v -race -covermode atomic -coverprofile $$(COVERAGE_DIR)/combined.txt -bench=. -benchmem -timeout 20m'

test-with-flags:
	@go test $(TEST_FLAGS) ./...
