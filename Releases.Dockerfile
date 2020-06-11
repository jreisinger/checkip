FROM golang:1.14-alpine AS build
RUN apk add --no-cache gcc libc-dev

WORKDIR /go/src/app
COPY . .
RUN go test ./...
RUN mkdir /releases

RUN GOOS=linux go build -o /releases/checkip
RUN tar -czvf /releases/checkip_linux.tar.gz -C /releases/ checkip
RUN rm /releases/checkip