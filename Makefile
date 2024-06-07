test:
	go test -cover ./...

install: test
	go install checkip.go

run: install
	checkip 91.228.166.47
