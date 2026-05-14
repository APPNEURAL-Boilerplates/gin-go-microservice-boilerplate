APP_NAME := gin-microservice

.PHONY: dev build test fmt fmt-check tidy check docker-up docker-down

dev:
	go run ./cmd/api

build:
	go build -o bin/$(APP_NAME) ./cmd/api

test:
	go test ./...

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './vendor/*')

fmt-check:
	@test -z "$$(gofmt -l $$(find . -name '*.go' -not -path './vendor/*'))" || \
		(echo "Go files need formatting:" && gofmt -l $$(find . -name '*.go' -not -path './vendor/*') && exit 1)

tidy:
	go mod tidy

check: fmt-check test

docker-up:
	docker compose up --build

docker-down:
	docker compose down
