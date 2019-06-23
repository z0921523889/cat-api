ARG GO_VERSION=1.12

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*

RUN mkdir -p /app
WORKDIR /app

COPY . .

RUN go get github.com/gin-gonic/gin

RUN go get github.com/go-xorm/xorm

RUN go get github.com/lib/pq

RUN github.com/gin-contrib/sessions

RUN go build -o ./app ./run/main.go

FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

RUN mkdir -p /app
WORKDIR /app
COPY --from=builder /app .

EXPOSE 8080

ENTRYPOINT ["./app"]