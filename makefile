APP_NAME=api
CMD_PATH=./cmd/api

.PHONY: help
help:
	@echo "Targets:"
	@echo "  make run        - run the API locally"
	@echo "  make test       - run unit tests"
	@echo "  make fmt        - format code"
	@echo "  make tidy       - go mod tidy"
	@echo "  make lint       - run golangci-lint (requires installation)"
	@echo "  make swagger    - generate swagger docs (added later)"
	@echo "  make build      - build binary to ./bin"

.PHONY: run
run:
	go run $(CMD_PATH)

.PHONY: test
test:
	go test ./... -race -count=1

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: build
build:
	mkdir -p bin
	go build -o bin/$(APP_NAME) $(CMD_PATH)

.PHONY: lint
lint:
	golangci-lint run ./...
