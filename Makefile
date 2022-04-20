test:
	go test ./...

install: test
	go install cmd/checkip/checkip.go

run: install
	checkip 1.1.1.1
	checkip -j 1.1.1.1 | jq -r \
'.checks[] | select(.type==1 or .type==2) | "\(.malicious) \(.name)"'
