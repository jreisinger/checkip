install:
	go install cmd/checkip.go

build:
	go build -ldflags "-w" -o cmd/checkip cmd/checkip.go

run: build
	./cmd/checkip 140.82.114.4
	./cmd/checkip 218.92.0.158
	./cmd/checkip 92.118.160.17

PLATFORMS := linux/amd64 darwin/amd64 linux/arm windows/amd64

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

release: test $(PLATFORMS)

$(PLATFORMS):
	GOOS=$(os) GOARCH=$(arch) go build -ldflags "-w" -o cmd/checkip-$(os)-$(arch) cmd/checkip.go
	tar -cf - cmd/checkip-$(os)-$(arch) | gzip -9c > cmd/checkip-$(os)-$(arch).tar.gz
	rm -f cmd/checkip-$(os)-$(arch)