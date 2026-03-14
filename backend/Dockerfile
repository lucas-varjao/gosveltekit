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
    ./main.go

FROM alpine:latest

WORKDIR /app

ARG APP_VERSION=dev

LABEL org.opencontainers.image.version=$APP_VERSION

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S -G app app

COPY --from=builder /out/server ./server
COPY configs ./configs

ENV TZ=America/Sao_Paulo
ENV GIN_MODE=release

EXPOSE 8080

USER app

CMD ["./server"]
