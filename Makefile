install:
	go install cmd/checkip.go

run: install
	checkip 140.82.114.4
	checkip 218.92.0.158
	checkip 92.118.160.17
	checkip -j 92.118.160.17 | jq -r '.[] | select(.Type=="Sec") | "\(.Name) => \(.IsMalicious)"'

PLATFORMS := linux/amd64 darwin/amd64 linux/arm windows/amd64

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

release: $(PLATFORMS)

$(PLATFORMS):
	GOOS=$(os) GOARCH=$(arch) go build -ldflags "-w" -o cmd/checkip-$(os)-$(arch) cmd/checkip.go
	tar -cf - cmd/checkip-$(os)-$(arch) | gzip -9c > cmd/checkip-$(os)-$(arch).tar.gz
	rm -f cmd/checkip-$(os)-$(arch)