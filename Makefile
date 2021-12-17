.PHONY: build
build:
	go build -v ./cmd/apiserver

run:
	go run ./cmd/apiserver


run:
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := build