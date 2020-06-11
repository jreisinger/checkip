test:
	go clean -testcache && go test -cover ./...

build: test
	go build

install: test
	go install
