# checkip

Checkip is a CLI tool and [library](https://pkg.go.dev/github.com/jreisinger/checkip) that provides generic and security information about IP addresses in an easy way.

```
$ checkip 1.1.1.1
db-ip.com       Sydney, Australia
dns name        one.one.one.one
iptoasn.com     CLOUDFLARENET
ping            0% packet loss (5/5), avg round-trip 15 ms
tls             TLS 1.3, exp. 2022/10/25, cloudflare-dns.com, *.cloudflare-dns.com, one.one.one.one
malicious       0% (0/7) ✅
```

## Usage examples

Select Sec (1) and InfoSec (2) check [type](https://pkg.go.dev/github.com/jreisinger/checkip#Type) and show which [check](https://pkg.go.dev/github.com/jreisinger/checkip/check) considers the IP address to be malicious:

```
$ checkip -j 1.1.1.1 | \
jq -r '.checks[] | select(.type==1 or .type==2) | "\(.malicious) \(.name)"'
false cinsscore.com
false firehol.org
false blocklist.de
false tls
false github.com/stamparm/ipsum
false threatcrowd.org
false otx.alienvault.com
```

Generate two random IP addresses and see if they are considered malicious:

```
$ ./randip 2 | checkip -a -j 2> /dev/null | \
jq -r '"\(.malicious_prob)\t\(.ipaddr)"'
0	176.214.10.86
0.1	229.236.76.24
```

See Wiki for more [usage examples](https://github.com/jreisinger/checkip/wiki/Usage-examples).

## Installation

To install the CLI tool

```
# optional; to install inside a container
docker run --rm -it golang /bin/bash

go install github.com/jreisinger/checkip/cmd/checkip@latest
```

or download a [release](https://github.com/jreisinger/checkip/releases) binary (from under "Assets") for your system and architecture.

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

## Development

Checkip is easy to extend. If you want to add a new way of checking an IP address, just write a function of type [Check](https://pkg.go.dev/github.com/jreisinger/checkip#Check). Add the new check to `check.All` [variable](https://pkg.go.dev/github.com/jreisinger/checkip/check#pkg-variables) and consider adding it to `check.Default` variable.

```
make run # test, install and run

git commit -m "improve tag docs" main.go

git tag | sort -V # or git ll
git tag -a v0.16.2 -m "improve docs"

git push --follow-tags
```
