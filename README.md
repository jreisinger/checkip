[![Go Reference](https://pkg.go.dev/badge/github.com/jreisinger/checkip.svg)](https://pkg.go.dev/github.com/jreisinger/checkip)
[![Go Report Card](https://goreportcard.com/badge/github.com/jreisinger/checkip)](https://goreportcard.com/report/github.com/jreisinger/checkip)

# checkip

Checkip is a CLI tool and library that provides generic and security information about IP addresses in an easy way.

```
$ checkip 91.228.166.47
db-ip.com       Petržalka, Slovakia
dns name        skh1-webredir01-v.eset.com
iptoasn.com     ESET-AS
ping            0% packet loss (5/5), avg round-trip 12 ms
shodan.io       OS: n/a, open: tcp/80 (nginx), tcp/443 (nginx), vulns: n/a
tls             TLS 1.3, exp. 2023/01/02, www.eset.com, eset.com
malicious       0% (0/8) ✅
```

## Usage examples

Check multiple IP addresses coming from stdin:

```
$ dig +short eset.sk | checkip
--- 91.228.167.128 ---
db-ip.com       Petržalka, Slovakia
dns name        h3-webredir02-v.eset.com
iptoasn.com     ESET-AS
ping            0% packet loss (5/5), avg round-trip 42 ms
shodan.io       OS: n/a, open: tcp/80 (nginx), tcp/443 (nginx), vulns: n/a
tls             TLS 1.3, exp. 2023/01/02, www.eset.com, eset.com
malicious       0% (0/8) ✅
--- 91.228.166.47 ---
db-ip.com       Petržalka, Slovakia
dns name        skh1-webredir01-v.eset.com
iptoasn.com     ESET-AS
ping            0% packet loss (5/5), avg round-trip 11 ms
shodan.io       OS: n/a, open: tcp/80 (nginx), tcp/443 (nginx), vulns: n/a
tls             TLS 1.3, exp. 2023/01/02, www.eset.com, eset.com
malicious       0% (0/8) ✅
```

Select Sec (1) and InfoSec (2) check [type](https://pkg.go.dev/github.com/jreisinger/checkip#Type) and show which [check](https://pkg.go.dev/github.com/jreisinger/checkip/check) considers the IP address to be malicious:

```
checkip -j 91.228.166.47 | \
jq -r '.checks[] | select(.type==1 or .type==2) | "\(.malicious) \(.name)"'
false firehol.org
false cinsscore.com
false tls
false blocklist.de
false github.com/stamparm/ipsum
false threatcrowd.org
false shodan.io
false otx.alienvault.com
```

Generate two random IP addresses and see if they are considered malicious:

```
$ ./randip 2 | checkip -a -j 2> /dev/null | \
jq -r '"\(.malicious_prob) \(.ipaddr)"'
0 53.18.151.128
0 163.201.51.56
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
