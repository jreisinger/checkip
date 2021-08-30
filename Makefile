test:
	go test -race -cover ./...

install: test
	go install cmd/checkip.go

run: install
	# ok ipaddr
	checkip 1.1.1.1 | sort
	checkip 8.8.8.8 | sort
	# suspicious ipaddr
	checkip 218.92.0.158 | sort
	checkip 92.118.160.17 | sort

build: test
	GOOS=linux GOARCH=amd64 go build -o cmd/checkip-linux-amd64 cmd/checkip.go
	GOOS=darwin GOARCH=amd64 go build -o cmd/checkip-darwin-amd64 cmd/checkip.go