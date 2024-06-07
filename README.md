[![Go Reference](https://pkg.go.dev/badge/github.com/jreisinger/checkip.svg)](https://pkg.go.dev/github.com/jreisinger/checkip)
[![Go Report Card](https://goreportcard.com/badge/github.com/jreisinger/checkip)](https://goreportcard.com/report/github.com/jreisinger/checkip)
[![StandWithUkraine](https://raw.githubusercontent.com/vshymanskyy/StandWithUkraine/main/badges/StandWithUkraine.svg)](https://github.com/vshymanskyy/StandWithUkraine/blob/main/docs/README.md)

# checkip

Sometimes I come across an IP address, for example when reviewing logs. And I'd like to find out more about this numerical label. Checkip is CLI tool and Go [library](https://pkg.go.dev/github.com/jreisinger/checkip/check) that provides infomation on and security posture of IP addresses. Most checks are passive and active checks are not aggressive.

## Usage examples

Check an IP address:

```
❯ checkip 91.228.166.47
--- 91.228.166.47 ---
db-ip.com       Petržalka, Slovakia
dns name        skh1-webredir01-v.eset.com
iptoasn.com     ESET-AS
is on AWS       false
ping            100% packet loss (5/0), avg round-trip 0 ms
tls             TLS 1.3, exp. 2024/01/02!!, www.eset.com, eset.com
malicious       14% (1/7) ✅
```

Check multiple IP addresses coming from STDIN:

```
❯ dig +short eset.sk | checkip
--- 91.228.167.128 ---
db-ip.com       Petržalka, Slovakia
dns name        h3-webredir02-v.eset.com
iptoasn.com     ESET-AS
is on AWS       false
ping            100% packet loss (5/0), avg round-trip 0 ms
tls             TLS 1.3, exp. 2024/01/02!!, www.eset.com, eset.com
malicious       14% (1/7) ✅
--- 91.228.166.47 ---
db-ip.com       Petržalka, Slovakia
dns name        skh1-webredir01-v.eset.com
iptoasn.com     ESET-AS
is on AWS       false
ping            100% packet loss (5/0), avg round-trip 0 ms
tls             TLS 1.3, exp. 2024/01/02!!, www.eset.com, eset.com
malicious       14% (1/7) ✅
```

Use more detailed JSON output to filter out those checks that consider the IP address to be malicious:

```
❯ checkip -j 91.228.166.47 | jq '.checks[] | select(.malicious == true)'
{
  "name": "tls",
  "type": "infoAndSecurity",
  "malicious": true,
  "info": {
    "SAN": [
      "www.eset.com",
      "eset.com"
    ],
    "Version": 772,
    "Expiry": "2024-01-02T23:59:59Z"
  }
}
```

Continuously generate [random IP addresses](https://github.com/jreisinger/checkip/blob/master/randip) and check them (hit Ctrl-C to stop):

```
❯ while true; do ./randip; sleep 2; done | checkip 2> /dev/null
--- 120.0.40.221 ---
db-ip.com       Zhoutou, China
iptoasn.com     CHINA169-BACKBONE CHINA UNICOM China169 Backbone
is on AWS       false
ping            100% packet loss (5/0), avg round-trip 0 ms
malicious       0% (0/6) ✅
--- 109.31.58.28 ---
db-ip.com       Paris, France
dns name        28.58.31.109.rev.sfr.net
iptoasn.com     LDCOMNET
is on AWS       false
ping            100% packet loss (5/0), avg round-trip 0 ms
malicious       0% (0/6) ✅
```

Generate 100 random IP addresses and select Russian or Chinese:

```
❯ ./randip 100 | checkip -c 20 -j 2> /dev/null | \
jq -r '.ipaddr as $ip | .checks[] | select (.name == "db-ip.com" and (.info.iso_code == "RU" or .info.iso_code == "CN")) | $ip'
43.33.161.208
223.13.196.8
```

Find out who is trying to SSH into your Linux system:

```
❯ sudo journalctl --unit ssh --since "1 hour ago" | \
∙ grep 'Bye Bye' | perl -wlne '/from ([\d\.]+)/ && print $1' | sort | uniq | \
∙ checkip 2> /dev/null
--- 167.172.105.64 ---
db-ip.com       Frankfurt am Main, Germany
iptoasn.com     DIGITALOCEAN-ASN
ping            0% packet loss (5/5), avg round-trip 21 ms
tls             TLS 1.3, exp. 2024/12/27, portal.itruck.com.sa, www.portal.itruck.com.sa
malicious       43% (3/7) 🤏
--- 180.168.95.234 ---
db-ip.com       Shanghai, China
iptoasn.com     CHINANET-SH-AP China Telecom Group
ping            0% packet loss (5/5), avg round-trip 213 ms
malicious       50% (3/6) 🚫
```

## Installation

To install the CLI tool

```
# optional; to install inside a container
docker run --rm -it golang /bin/bash

go install github.com/jreisinger/checkip/cmd/checkip@latest
```

or download a [release](https://github.com/jreisinger/checkip/releases) binary (from under "Assets") for your system and architecture.

## Configuration

For some checks to start working you need to register and get an API (LICENSE) key. See the service web site for how to do that. An absent key is not reported as an error, the check is simply not executed and `missing_credentials` JSON field is set.

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
