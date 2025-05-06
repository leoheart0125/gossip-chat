GOBIN?=$(shell go env GOPATH)/bin
GOENTRY?=./cmd/chat/main.go

.PHONY: install
install:
	@go mod tidy
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest


.PHONY: lint
lint:
	@go fmt ./...
	@$(GOBIN)/golangci-lint run --fix

.PHONY: run
run:
	@go run $(GOENTRY)