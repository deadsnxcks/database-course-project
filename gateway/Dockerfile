FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /out/gateway ./gateway/cmd/dbcp-api-gateway

FROM alpine:3.20

WORKDIR /app

RUN adduser -D -u 10001 app && chown -R app:app /app

COPY --from=builder --chown=app:app /out/gateway .
COPY --from=builder --chown=app:app /src/gateway/config ./config

USER app

ENV CONFIG_PATH=/app/config/docker.yaml

EXPOSE 8081

CMD ["./gateway"]