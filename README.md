# checkip

Checkip is a CLI tool and library that provides generic and security information about an IP address in an easy way. It uses various [checks](https://pkg.go.dev/github.com/jreisinger/checkip/check) to do so.

```
$ checkip 1.1.1.1
db-ip.com   Sydney, Australia
dns name    one.one.one.one
iptoasn.com CLOUDFLARENET
ping        0% packet loss, sent 5, recv 5, avg round-trip 14 ms
tls         TLS 1.3, cloudflare-dns.com, *.cloudflare-dns.com, one.one.one.one
malicious   0% (0/6) ✅
```

You can get output also in JSON (`-j`). Here we select Sec (1) and InfoSec (2) check [type](https://pkg.go.dev/github.com/jreisinger/checkip#Type) and show if the check considers the IP address to be malicious.

```
$ checkip -j 1.1.1.1 | jq -r \
'.checks[] | select(.type==1 or .type==2) | "\(.malicious)\t\(.name)"'
false	blocklist.de
false	cinsscore.com
false	firehol.org
false	github.com/stamparm/ipsum
false	otx.alienvault.com
false	threatcrowd.org
```

See [wiki](https://github.com/jreisinger/checkip/wiki) for more.

## Installation

To install the CLI tool

```
# optional; to install inside a container
docker run --rm -it golang /bin/bash

go install github.com/jreisinger/checkip/cmd/checkip@latest
```

or download a [release](https://github.com/jreisinger/checkip/releases) binary (from under "Assets") for your system and architecture.

## Development

Checkip is easy to extend. If you want to add a new way to check an IP address, just write a function of type [Check](https://pkg.go.dev/github.com/jreisinger/checkip#Check). Consider adding the check to `check.Default` [variable](https://pkg.go.dev/github.com/jreisinger/checkip/check#pkg-variables).

```
make run # test, install and run

git commit -m "improve tag docs" main.go

git tag | sort -V # or git ll
git tag -a v0.16.2 -m "improve docs"

git push --follow-tags
```

## Configuration

For some checks to work you need to register and get an API (LICENSE) key. See the service web site for how to do that.

Store the keys in `$HOME/.checkip.yaml` file.

```
ABUSEIPDB_API_KEY: aaaaaaaabbbbbbbbccccccccddddddddeeeeeeeeffffffff11111111222222223333333344444444
MAXMIND_LICENSE_KEY: abcdef1234567890
SHODAN_API_KEY: aaaabbbbccccddddeeeeffff11112222
URLSCAN_API_KEY: abcd1234-a123-4567-678z-a2b3c4b5d6e7
VIRUSTOTAL_API_KEY: aaaaaaaabbbbbbbbccccccccddddddddeeeeeeeeffffffff1111111122222222
```

You can also use environment variables with the same names.