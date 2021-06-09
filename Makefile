test:
	go test ./...

install: test
	go install cmd/checkip.go

run: install
	checkip 1.1.1.1 # ok
	checkip 218.92.0.158 # suspicious