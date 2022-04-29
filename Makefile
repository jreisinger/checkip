test:
	go test -cover ./...

install: test
	go install cmd/checkip/checkip.go

run: install
	checkip 91.228.166.47
	dig +short eset.sk | checkip
	checkip -j 91.228.166.47 | \
jq -r '.checks[] | select(.type==1 or .type==2) | "\(.malicious) \(.name)"'
	./randip 2 | checkip -a -j 2> /dev/null | \
jq -r '"\(.malicious_prob) \(.ipaddr)"'
