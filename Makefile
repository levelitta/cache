LOCAL_BIN:=$(CURDIR)/bin
GOBIN=$(LOCAL_BIN)

.PHONY: .bin-deps
.bin-deps:
	GOBIN=$(LOCAL_BIN) go install golang.org/x/tools/cmd/benchcmp

.PHONY: generate
generate: .bin-deps .generate

.PHONY: generate-fast
generate-fast: .generate

.PHONY: .generate
.generate:
	go generate ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: test-race
test-race:
	go test --race -v ./...

.PHONY: bench
bench:
	go test -benchmem -bench=. ./... -run=^# | tee ./benchmarks/results/new.txt