test:
	go test -race -cover ./...

install: test
	go install

run: install
	checkip 91.228.166.47
	checkip 218.92.0.158
	checkip -a 45.33.32.156 # scanme.nmap.org
