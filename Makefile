test:
	go test -race -cover ./...

install: test
	go install cmd/checkip.go

run: install
	checkip 140.82.114.4 | sort
	checkip 218.92.0.158 | sort
	checkip 92.118.160.17 | sort

PLATFORMS := linux/amd64 darwin/amd64 linux/arm windows/amd64

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

release: test $(PLATFORMS)

$(PLATFORMS):
	GO111MODULE=on GOOS=$(os) GOARCH=$(arch) go build -o cmd/checkip-$(os)-$(arch) cmd/checkip.go