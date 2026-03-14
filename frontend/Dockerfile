# syntax=docker/dockerfile:1

FROM oven/bun:1.2.23-alpine AS builder

WORKDIR /app

ARG APP_VERSION=dev
ARG VITE_API_URL=http://localhost:8080
ENV APP_VERSION=${APP_VERSION}
ENV VITE_API_URL=${VITE_API_URL}

COPY package.json bun.lock ./
RUN bun install --frozen-lockfile

COPY . .

RUN bun run build

FROM nginx:1.29-alpine

ARG APP_VERSION=dev

LABEL org.opencontainers.image.version=$APP_VERSION

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/build /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
