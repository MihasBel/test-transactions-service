.PHONY: mod
mod:
	go mod download

.PHONY: build
build:
	go build \
			--trimpath \
			-o bin/app/product-details \
			cmd/http/main.go

.PHONY: check
check:
	golangci-lint run -v --config .golangci.yml

.PHONY: test
test:
	go test -v ./...

.PHONY: race
race:
	go test -v -race ./...

.PHONY: run
run:
	go run cmd/http/main.go --config=configs/local/env.json

.PHONY: swag
swag:
	swag init -g cmd/http/main.go

.PHONY: generate
generate:
	@go generate ./...