VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

.PHONY: build run test clean sqlc migrate-up seed

build:
	go build -ldflags "-X main.version=$(VERSION)" -o bin/inkvoice ./cmd/inkvoice

run: build
	./bin/inkvoice

test:
	go test ./...

clean:
	rm -rf bin/ inkvoice.db

sqlc:
	sqlc generate

migrate-up: build
	./bin/inkvoice migrate up

seed: build
	./bin/inkvoice seed template
	./bin/inkvoice seed data --skip-if-exists
