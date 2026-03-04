GO ?= go
DEMO_PKG := ./cmd/eip4844-demo

.PHONY: help fmt test run examples check

help:
	@echo "Available targets:"
	@echo "  make fmt       - format Go code"
	@echo "  make test      - run all tests"
	@echo "  make run       - run the demo app"
	@echo "  make examples  - run example cases"
	@echo "  make check     - fmt + test"

fmt:
	$(GO) fmt ./...

test:
	$(GO) test ./...

run:
	$(GO) run $(DEMO_PKG)

examples: run

check: fmt test

#make
