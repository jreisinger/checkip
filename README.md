[![Go Reference](https://pkg.go.dev/badge/github.com/jreisinger/checkip.svg)](https://pkg.go.dev/github.com/jreisinger/checkip)
[![Go Report Card](https://goreportcard.com/badge/github.com/jreisinger/checkip)](https://goreportcard.com/report/github.com/jreisinger/checkip)
[![StandWithUkraine](https://raw.githubusercontent.com/vshymanskyy/StandWithUkraine/main/badges/StandWithUkraine.svg)](https://github.com/vshymanskyy/StandWithUkraine/blob/main/docs/README.md)

# checkip

Sometimes I come across an IP address, reviewing logs for example, and I want to know more about this numerical label. Checkip is CLI tool and Go [library](https://pkg.go.dev/github.com/jreisinger/checkip/check) that provides (security) information on IP addresses. It runs various checks to get the information. Most checks are passive, i.e. not interacting directly with the IP address. Active checks, like ping and tls, are not aggressive and can be disabled with `-no-active`.

## Quick start

```
$ go install github.com/jreisinger/checkip@latest
$ checkip 1.1.1.1 91.228.166.47
--- 1.1.1.1 ---
db-ip.com       Sydney, Australia
dns name        one.one.one.one
iptoasn.com     CLOUDFLARENET
is on AWS       false
isc.sans.edu    attacks: 0, abuse contact: abuse@cloudflare.com
ping            0% packet loss (5/5), avg round-trip 4 ms
tls             TLS 1.3, exp. 2026/01/21, cloudflare-dns.com, *.cloudflare-dns.com, one.one.one.one
virustotal.com  network: 1.1.1.0/24, SAN: cloudflare-dns.com, *.cloudflare-dns.com, one.one.one.one, 1.0.0.1, 1.1.1.1, 162.159.36.1, 162.159.46.1, 2606:4700:4700::1001, 2606:4700:4700::1111, 2606:4700:4700::64, 2606:4700:4700::6400
malicious prob. 0% (0/12) ✅
--- 91.228.166.47 ---
db-ip.com       Petržalka, Slovakia
dns name        skh1-webredir01-v.eset.com
iptoasn.com     ESET-AS
is on AWS       false
isc.sans.edu    attacks: 0, abuse contact: domains@eset.sk
ping            0% packet loss (5/5), avg round-trip 5 ms
tls             TLS 1.3, exp. 2024/01/02!!, www.eset.com, eset.com
virustotal.com  network: 91.228.164.0/22, SAN: www.eset.com, eset.com
malicious prob. 17% (2/12) 🤏
```

## More usage examples

Use detailed JSON output to filter out those checks that consider the IP address to be malicious:

```
checkip -j 91.228.166.47 | jq '.checks[] | select(.ipAddrIsMalicious == true)'
```

Check multiple IP addresses coming from STDIN:

```
dig +short eset.sk | checkip
```

Continuously generate [random IP addresses](https://github.com/jreisinger/checkip/blob/master/randip) and check them (hit Ctrl-C to stop):

```
while true; do ./randip; sleep 2; done | checkip 2> /dev/null
```

Generate 100 random IP addresses and select Russian or Chinese:

```
./randip 100 | checkip -p 20 -j 2> /dev/null | \
jq -r '.ipAddr as $ip | .checks[] | select (.description == "db-ip.com" and (.ipAddrInfo.iso_code == "RU" or .ipAddrInfo.iso_code == "CN")) | $ip'
```

Find out who is trying to SSH into your Linux system:

```
sudo journalctl --unit ssh --since "1 hour ago" | \
grep 'Bye Bye' | perl -wlne '/from ([\d\.]+)/ && print $1' | sort | uniq | \
checkip 2> /dev/null
```

## Installation

To install the CLI tool

```
# optional; to install inside a container
docker run --rm -it golang /bin/bash

go install github.com/jreisinger/checkip@latest
```

or download a [release](https://github.com/jreisinger/checkip/releases) binary (from under "Assets") for your system and architecture.

## Configuration and cache

For some checks to start working you need to sign up on a web site (like https://www.abuseipdb.com) and get an API (or LICENSE) key. Checkip doesn't report an absent API key as an error; the check is simply not executed and `missingCredentials` JSON field is set to the name of the API key (like `ABUSEIPDB_API_KEY`).

Store the keys in `$HOME/.checkip.yaml` file:

```
ABUSEIPDB_API_KEY: aaaaaaaabbbbbbbbccccccccddddddddeeeeeeeeffffffff11111111222222223333333344444444
MAXMIND_LICENSE_KEY: abcdef1234567890
SHODAN_API_KEY: aaaabbbbccccddddeeeeffff11112222
URLSCAN_API_KEY: abcd1234-a123-4567-678z-a2b3c4b5d6e7
VIRUSTOTAL_API_KEY: aaaaaaaabbbbbbbbccccccccddddddddeeeeeeeeffffffff1111111122222222
```

You can also use environment variables with the same names.

Data used by some checks is downloaded (cached) to `$HOME/.checkip/` folder. Is gets periodically re-downloaded so it is fresh.

Repeated checks of the same IP address within a single `checkip` process are memoized (cached) in memory, so duplicate inputs in one run reuse already computed results instead of hitting the same providers again.

Some passive remote checks also persist per-IP results under `$HOME/.checkip/results/v1/` for one hour, so repeated runs can reuse recent results without hitting the same APIs again. Use `checkip -no-cache ...` to bypass the in-memory and persistent result caches for a run. This does not bypass the source-data cache for downloaded files.

## Development

Checkip is easy to extend. If you want to add a new way of checking IP addresses:

1. Write a function of type [check.Func](https://pkg.go.dev/github.com/jreisinger/checkip/check#Func).
2. Add a [check.Definition](https://pkg.go.dev/github.com/jreisinger/checkip/check#Definition) to [check.Definitions](https://pkg.go.dev/github.com/jreisinger/checkip/check#Definitions).

Checks use per run (process-lifetime) memoization by default. If a check should persist results across runs, set `PersistentTTL` and provide `NewInfo` for decoding cached `ipAddrInfo` JSON. If a check directly contacts the target IP address, mark it as active so CLI users can disable it with `-no-active`. If a check must always run live, set its cache policy to `check.CacheNone`.

Typical workflow:

```
make run # test and run

git commit

git tag | sort -V | tail -1
git tag -a v0.2.0 -m "new check func"

git push --follow-tags # will build a new release on GitHub
```


## Extended branch

With a plugin system to choose each check on command line, or in config file, to manage your subscriptions plans

add MISP event output and experimental cotation scale

```
❯ git checkout extend
❯ make extend

❯ checkipext -h
Usage of checkipext:
 checkipext [-flag] IP [IP liste]
  -a string
        append to list of checks
  -d	debug
  -j	detailed output in JSON
  -m	MISP event output in JSON
  -p n
        check n IP addresses in parallel (default 5)
  -t string
        list of checks

  Available Checks :
  IOCLoc, AbuseIPDB, IsOnAWS, OTX, Tls, UrlScan, DnsMX, DnsName, Misp, MyDB, Onyphe, SansISC, BlockList, DBip, IPSum,
  IPtoASN, MaxMind, Spur, CinsScore, Censys, Firehol, IpAPI, Ping, Shodan, VirusTotal
```

or add default checks in  ``$HOME/.checkip.yaml`` file
```
CHECKS: IOCLoc, IpAPI, MyDB, Spur, BlockList, CinsScore, DBip, DnsName, Firehol, IPSum, IPtoASN, IsOnAWS, OTX, AbuseIPDB, Shodan, Onyphe, Tls
```

```
❯ checkipext 91.228.166.47
Checks: IOCLoc,IpAPI,MyDB,Spur,BlockList,CinsScore,DBip,DnsName,Firehol,IPSum,IPtoASN,IsOnAWS,OTX,
 AbuseIPDB,Shodan,Onyphe,Tls
--- 91.228.166.47 ---
IOCLoc          91.228.166.47 (SK)🇸🇰  AS50881 - ESET, spol. s r.o.
abuseipdb.com   domain: eset.com, usage type: Commercial
db-ip.com       Petržalka, Slovakia
dns name        skh1-webredir01-v.eset.com
iptoasn.com     ESET-AS
is on AWS       false
shodan.io       OS: n/a, open: tcp/80 (nginx), tcp/443 (nginx), vulns: n/a
tls             TLS 1.3, exp. 2024/01/02!!, www.eset.com, eset.com
Cotation        A1 - server
malicious prob. 11% (1/9) ✅

--- 148.72.164.179 ---
IOCLoc          148.72.164.179 (US)🇺🇸  AS30083 - AS-30083-US-VELIA-NET
abuseipdb.com   domain: velia.net, usage type: Data Center/Web Hosting/Transit
ipapi.is        VPN (NordVPN)
db-ip.com       St Louis, United States
iptoasn.com     AS-30083-US-VELIA-NET
spur.io         VPN : NORD_VPN
is on AWS       false
Cotation        B1 - vpn
malicious prob. 12% (1/8) ✅

IOC: 91.228.166.47 (SK)🇸🇰  AS50881 - ESET, spol. s r.o. [A1 - server],
   148.72.164.179 (US)🇺🇸  AS30083 - AS-30083-US-VELIA-NET [B1 - vpn]
```

With additional and optional checks :
  - Onyphe: https://www.onyphe.io/
  - IpAPI: https://ipapi.is
  - IOCLoc : list all "IP (country) ASN"
  - MyDB : to check your own DB by IP ``curl -H "Authorization: bearer {{token}} "MYDB_URL/{{IP}}"``
  - Misp : to check by attribute ip-src on your own Misp instance


add keys in ``$HOME/.checkip.yaml`` file

```
MISP_URL: https://localhost
MISP_KEY: xxxxxxxxxxxxxxxxxxxxxxxx
# optional MISP_OPT :
#  selfsigned cert, search in last 365 days (or h for hours)
MISP_OPT: selfsigned, 365d

MYDB_URL: https://zzzzzzzzzz/sss/ssss
MYDB_API_KEY: xxxxxxxxxxx

ONYPHE_API_KEY: xxxxxxxxxxxxxxx

IP_API_KEY : xxxxxxxxxxxxxxx
```
