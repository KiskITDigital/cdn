ARG GO_VER
ARG ALPINE_VER

FROM golang:1.23.1-alpine3.20 as builder

WORKDIR /src
COPY . .

RUN apk upgrade --no-cache && \
    CGO_ENABLED=0 GOOS=linux go build -mod=readonly -a -installsuffix cgo -o cdn ./cmd/cdn

FROM alpine:3.20

RUN addgroup -g 101 app && \
    adduser -H -u 101 -G app -s /bin/sh -D app && \
    apk update --no-cache && \
    apk upgrade --no-cache

WORKDIR /app
COPY --from=builder --chown=app:app /src/cdn .
COPY --from=builder --chown=app:app /src/swaggerui/ ./swaggerui/

USER app

CMD ["/app/cdn"]
