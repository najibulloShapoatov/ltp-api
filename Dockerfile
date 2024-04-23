# Build stage
FROM golang:1.22-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN apk add --no-cache git && \
    go build -o ./bin/app cmd/main.go

# Run stage
FROM alpine:3.18.6
WORKDIR /app
COPY --from=builder /app/bin/app .
COPY --from=builder /app/docs docs
RUN ln -s /etc/config.yaml /app/config.yaml
RUN apk update && apk add --no-cache ca-certificates && rm -rf /var/cache/apk/*
ENTRYPOINT ["./app"]