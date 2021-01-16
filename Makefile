VERSION ?= dev

test:
	go clean -testcache && go test -race -cover ./...

install: test
	go install -ldflags "-X main.Version=${VERSION}"

release:
	docker build --build-arg version=${VERSION} -t checkip-releases -f Releases.Dockerfile .
	docker create -ti --name checkip-releases checkip-releases sh
	test -d releases || mkdir releases
	docker cp checkip-releases:/releases/checkip_linux_amd64 releases/
	docker cp checkip-releases:/releases/checkip_darwin_amd64 releases/
	docker rm checkip-releases
	docker rmi checkip-releases