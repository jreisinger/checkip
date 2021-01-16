ARG VERSION=dev
FROM golang:1.15

WORKDIR /go/src/app
COPY . .
RUN go test ./...
RUN mkdir /releases

ARG version=dev
RUN GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$version" -o /releases/checkip-linux-amd64
RUN GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.Version=$version" -o /releases/checkip-darwin-amd64