test:
	go test -v -race -cover ./...
	go vet

install: test
	go install

run: install
	checkip 91.228.166.47
	checkip 209.141.33.65
	checkip 218.92.0.158
