#build stage
FROM golang:1.21-alpine3.17 AS builder

WORKDIR /app/

COPY . .

RUN go mod download
RUN set -x; apk add --no-cache && CGO_ENABLED=0 go build -ldflags="-s -w" -o ./bin/app cmd/main.go


#Run stage
FROM alpine:3.17

WORKDIR /app

COPY --from=0 /app/bin .
COPY --from=0 /app/docs docs


RUN ln -s /etc/config.yaml /app/config.yaml
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

ENTRYPOINT ["./app"]
