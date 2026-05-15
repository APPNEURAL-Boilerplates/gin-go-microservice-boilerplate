# Gin Go Microservice

A production-friendly Go microservice starter using Gin.

## Stack

- Go 1.25+
- Gin 1.12
- Standard library structured logging with `log/slog`
- REST API with versioned routes under `/api/v1`
- Docker and Docker Compose
- GitHub Actions CI

## Project structure

```txt
gin-go-microservice/
├─ cmd/api/main.go
├─ internal/
│  ├─ app/
│  ├─ clients/
│  ├─ common/
│  ├─ config/
│  ├─ events/
│  ├─ middleware/
│  ├─ modules/
│  │  ├─ health/
│  │  ├─ items/
│  │  └─ root/
│  └─ workers/
├─ tests/
├─ Dockerfile
├─ docker-compose.yml
├─ Makefile
├─ go.mod
└─ README.md
```

## Features

- App/router factory pattern
- Graceful shutdown
- JSON request logging
- Request ID middleware with `X-Request-Id`
- Security headers middleware
- Panic recovery with JSON error response
- Consistent success/error response format
- Health and readiness endpoints
- Example items module
- Handler → service → repository structure
- In-memory repository placeholder
- Outbound HTTP client placeholder
- Event publisher placeholder
- Background worker placeholder
- Integration tests with `net/http/httptest`

## Run locally

```bash
cp .env.example .env
go mod tidy
go run ./cmd/api
```

Open:

```txt
http://localhost:8080
http://localhost:8080/api/v1/health
http://localhost:8080/api/v1/ready
```

## Main endpoints

```txt
GET  /                     Service metadata
GET  /api/v1/health        Health check
GET  /api/v1/ready         Readiness check
GET  /api/v1/items         List items
POST /api/v1/items         Create item
GET  /api/v1/items/:id     Get item by ID
```

## Example request

```bash
curl -X POST http://localhost:8080/api/v1/items \
  -H "content-type: application/json" \
  -H "x-request-id: demo-request-1" \
  -d '{"name":"Keyboard","description":"Mechanical keyboard","price":99.99}'
```

## Tests

```bash
go test ./...
```

Or:

```bash
make check
```

## Docker

```bash
cp .env.example .env
docker compose up --build
```

## Environment variables

```txt
SERVICE_NAME=gin-microservice
ENVIRONMENT=development
HOST=0.0.0.0
PORT=8080
GIN_MODE=debug
LOG_LEVEL=info
READ_TIMEOUT_SECONDS=10
WRITE_TIMEOUT_SECONDS=10
SHUTDOWN_TIMEOUT_SECONDS=10
TRUSTED_PROXIES=
```

## Response format

Success:

```json
{
  "ok": true,
  "data": {},
  "request_id": "req_..."
}
```

Error:

```json
{
  "ok": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "route not found"
  },
  "request_id": "req_..."
}
```

## Production notes

- Put real persistence behind the repository interface.
- Replace `LogPublisher` with Kafka, RabbitMQ, NATS, SQS, Pub/Sub, or your event broker.
- Keep secrets out of `.env.example` and source control.
- Set `GIN_MODE=release` in production.
- Configure `TRUSTED_PROXIES` when running behind a trusted reverse proxy.
- Add authentication/authorization middleware before exposing private endpoints.
