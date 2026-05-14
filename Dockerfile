FROM golang:1.26-alpine AS build

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /bin/service ./cmd/api

FROM alpine:3.22

RUN addgroup -S app && adduser -S app -G app && apk add --no-cache ca-certificates
USER app

COPY --from=build /bin/service /service

EXPOSE 8080
ENTRYPOINT ["/service"]
