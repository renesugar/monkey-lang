.PHONY: dev build install image release profile bench test clean

CGO_ENABLED=0
COMMIT=$(shell git rev-parse --short HEAD)

all: dev

dev: build
	@./monkey-lang -d

build: clean
	@go build \
		-tags "netgo static_build" -installsuffix netgo \
		-ldflags "-w -X $(shell go list)/version/.GitCommit=$(COMMIT)" \
		.

install: build
	@go install

image:
	@docker build -t prologic/monkey-lang .

release:
	@./tools/release.sh

profile:
	@go test -cpuprofile cpu.prof -memprofile mem.prof -v -bench ./...

bench:
	@go test -v -benchmem -bench=. ./...

test:
	@go test -v -cover -coverprofile=coverage.txt -covermode=atomic -coverpkg=./... -race ./...

clean:
	@git clean -f -d -X
