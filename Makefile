.PHONY: run
run:
	go run ./cmd/image-upload

.PHONY: build
build:
	go build -o server -v ./cmd/image-upload

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := run