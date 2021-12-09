test:
	go test -race -cover ./...

install: test
	go install

run: install
	checkip 91.228.166.47
	checkip 209.141.33.65 2> /dev/null
	checkip 218.92.0.158 2> /dev/null
