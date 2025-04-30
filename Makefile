GOBIN?=$(shell go env GOPATH)/bin

.PHONY: lint
lint:
	@go fmt ./...
	@$(GOBIN)/golangci-lint run --fix
