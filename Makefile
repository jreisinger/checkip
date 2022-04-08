test:
	go test ./...

install: test
	go install cmd/checkip/checkip.go

run: install
	checkip 218.92.0.158
	checkip -a 45.33.32.156 # scanme.nmap.org
	checkip -j 218.92.0.158 | jq -r \
'.checks[] | select((.type==1 or .type==2) and .info!=null) | "\(.malicious)\t\(.name)"'
