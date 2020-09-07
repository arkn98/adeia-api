.DEFAULT_GOAL := help

.PHONY: test test-coverage build help

# run all unit tests
test:
	go test -v -race ./...

# run all unit tests and report coverage data to codecov.io
test-coverage:
	go test -v -coverprofile=coverage.txt -covermode=atomic ./... && (curl -s https://codecov.io/bash | bash)

# build project
build:
	go build -v ./...

help:
	@echo "Usage:"
	@echo "make [command]"
	@echo
	@echo "Available commands:"
	@echo "build         : Build project"
	@echo "help          : Display this help message"
	@echo "test          : Run all unit-tests"
	@echo "test-coverage : Run all unit-tests and report coverage to codecov.io"
