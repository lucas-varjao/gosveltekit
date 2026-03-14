# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS builder

WORKDIR /src

ARG APP_VERSION=dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w -X 'gosveltekit/internal/version.BuildVersion=${APP_VERSION}'" \
    -o /out/server \
    .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w -X 'gosveltekit/internal/version.BuildVersion=${APP_VERSION}'" \
    -o /out/migrate \
    ./cmd/migrate

FROM alpine:3.22 AS runtime-base

WORKDIR /app

ARG APP_VERSION=dev

LABEL org.opencontainers.image.version=$APP_VERSION

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S -G app app

COPY configs ./configs

ENV TZ=America/Sao_Paulo
ENV GIN_MODE=release

USER app

FROM runtime-base AS backend

COPY --from=builder /out/server ./server

EXPOSE 8080

CMD ["./server"]

FROM runtime-base AS migrator

COPY --from=builder /out/migrate ./migrate
COPY db/migrations ./db/migrations

CMD ["./migrate", "up"]
