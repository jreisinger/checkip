FROM golang:1.14

WORKDIR /go/src/app
COPY . .
RUN go test ./...
RUN mkdir /releases

RUN GOOS=linux GOARCH=amd64 go build -o /releases/checkip
RUN tar -czvf /releases/checkip_linux_amd64.tar.gz -C /releases/ checkip
RUN rm /releases/checkip