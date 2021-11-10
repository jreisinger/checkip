test:
	go test -cover ./...

install: test
	go install 

run: install
	checkip 91.228.166.47
	checkip 209.141.33.65
	checkip 218.92.0.158

test-api:
	curl -s localhost:8000/api/v1/91.228.166.47 | jq -r '.Infos[] | select(.Info!="") | "\(.Name)\t\(.Info)"'
	curl -s localhost:8000/api/v1/209.141.33.65 | jq -r '.Infos[] | select(.Info!="") | "\(.Name)\t\(.Info)"'
	curl -s localhost:8000/api/v1/218.92.0.158 | jq -r '.Infos[] | select(.Info!="") | "\(.Name)\t\(.Info)"'
	curl -s 'localhost:8000/api/v1/91.228.166.47' | jq -r '.ProbabilityMalicious'
	curl -s 'localhost:8000/api/v1/209.141.33.65' | jq -r '.ProbabilityMalicious'
	curl -s 'localhost:8000/api/v1/218.92.0.158' | jq -r '.ProbabilityMalicious'
