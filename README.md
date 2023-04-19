[![Go Reference](https://pkg.go.dev/badge/github.com/jreisinger/checkip.svg)](https://pkg.go.dev/github.com/jreisinger/checkip)
[![Go Report Card](https://goreportcard.com/badge/github.com/jreisinger/checkip)](https://goreportcard.com/report/github.com/jreisinger/checkip)
[![StandWithUkraine](https://raw.githubusercontent.com/vshymanskyy/StandWithUkraine/main/badges/StandWithUkraine.svg)](https://github.com/vshymanskyy/StandWithUkraine/blob/main/docs/README.md)

# checkip

Sometimes I come across an IP address, for example when reviewing logs. And I'd like to find out more about this numerical label. Checkip is CLI tool and Go [library](https://pkg.go.dev/github.com/jreisinger/checkip/check) that provides generic and security information about IP addresses in a quick way.

```sh
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

```sh
$ checkip -j 34.250.182.30 | jq '.checks[] | select(.malicious == true)'
{
  "name": "shodan.io",
  "type": 2,
  "malicious": true,
  "info": {
    "org": "Amazon Data Services Ireland Limited",
    "data": [
      {
        "product": "lighttpd",
        "version": "1.4.53",
        "port": 80,
        "transport": "tcp"
      },
      {
        "product": "AWS ELB",
        "version": "2.0",
        "port": 443,
        "transport": "tcp"
      }
    ],
    "os": "",
    "ports": [
      80,
      443
    ],
    "vulns": [
      "CVE-2022-22707",
      "CVE-2019-11072"
    ]
  }
}
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

Data used by some checks are downloaded (cached) to `$HOME/.checkip/` folder. They are periodically re-downloaded so they are fresh.

## Development

Checkip is easy to extend. If you want to add a new way of checking IP addresses:

1. Write a function of type [Check](https://pkg.go.dev/github.com/jreisinger/checkip#Check). 
2. Add the new check to `check.All` [variable](https://pkg.go.dev/github.com/jreisinger/checkip/check#pkg-variables)
3. Optional: consider adding the new check to `check.Default` variable.

Typical workflow:

```
make run # test, install and run

git commit -m "backwards compatible bug fix" main.go

git tag | sort -V | tail -1
git tag -a v0.16.1 -m "patch" # will build a new release on GitHub when pushed

git push --follow-tags
```
