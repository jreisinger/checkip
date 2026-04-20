test:
	go test -race -cover ./...

run: test
	go run -race ./checkip.go 91.228.166.47

install: run
	go install ./checkip.go

extend: test
	go install ./cmd/extend/checkipext.go
