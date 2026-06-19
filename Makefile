.PHONY: build check lint test

ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))

lint test build: check

check:
	@"$(ROOT)/scripts/check-baseline.sh"
