test:
	go test ./...

run: test
	go run -race ./checkip.go 91.228.166.47

install: run
	go install ./checkip.go
