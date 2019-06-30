ARG GO_VERSION=1.12

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*

# Configure Go
ENV GOPATH /go
RUN mkdir -p $GOPATH/src/cat-api
WORKDIR $GOPATH/src/cat-api
COPY . .

RUN go get -u github.com/swaggo/swag/cmd/swag

RUN go get -d -v ./src/run/...

RUN swag init -d ./src/app/router -g router.go -o ./src/app/docs

RUN go build -o /app ./src/run/main.go

FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

COPY --from=builder /app .

EXPOSE 8085

ENTRYPOINT ["./app"]