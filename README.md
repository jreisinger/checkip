[![Go Reference](https://pkg.go.dev/badge/github.com/jreisinger/checkip.svg)](https://pkg.go.dev/github.com/jreisinger/checkip)
[![Go Report Card](https://goreportcard.com/badge/github.com/jreisinger/checkip)](https://goreportcard.com/report/github.com/jreisinger/checkip)
[![StandWithUkraine](https://raw.githubusercontent.com/vshymanskyy/StandWithUkraine/main/badges/StandWithUkraine.svg)](https://github.com/vshymanskyy/StandWithUkraine/blob/main/docs/README.md)

# checkip

Sometimes I come across an IP address, for example when reviewing logs. And I want to know something more about it. Checkip is CLI tool and Go library that provides generic and security information about IP addresses in a quick way.

```
$ checkip 91.228.166.47
--- 91.228.166.47 ---
db-ip.com       Petržalka, Slovakia
dns name        skh1-webredir01-v.eset.com
iptoasn.com     ESET-AS
ping            0% packet loss (5/5), avg round-trip 12 ms
shodan.io       OS: n/a, open: tcp/80 (nginx), tcp/443 (nginx), vulns: n/a
tls             TLS 1.3, exp. 2023/01/02, www.eset.com, eset.com
malicious       0% (0/8) ✅
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

NOTE: depending on your Internet connection the first run can take a while because data is downloaded (to `$HOME/.checkip`).

## Configuration

For some checks to start working you need to register and get an API (LICENSE) key. See the service web site for how to do that. An absent key is not reported as an error, the check is simply ignored.

Store the keys in `$HOME/.checkip.yaml` file:

```
ABUSEIPDB_API_KEY: aaaaaaaabbbbbbbbccccccccddddddddeeeeeeeeffffffff11111111222222223333333344444444
MAXMIND_LICENSE_KEY: abcdef1234567890
SHODAN_API_KEY: aaaabbbbccccddddeeeeffff11112222
URLSCAN_API_KEY: abcd1234-a123-4567-678z-a2b3c4b5d6e7
VIRUSTOTAL_API_KEY: aaaaaaaabbbbbbbbccccccccddddddddeeeeeeeeffffffff1111111122222222
```

You can also use environment variables with the same names.

## Development

Checkip is easy to extend. If you want to add a new way of checking IP addresses, just write a function of type [Check](https://pkg.go.dev/github.com/jreisinger/checkip#Check). Add the new check to `check.All` [variable](https://pkg.go.dev/github.com/jreisinger/checkip/check#pkg-variables) and consider adding it to `check.Default` variable.

```
make run # test, install and run

git commit -m "backwards compatible bug fix" main.go

git tag | sort -V               # or git ll
git tag -a v0.16.1 -m "patch"   # will build a new release

git push --follow-tags
```
