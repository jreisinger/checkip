# checkip

Checkip provides generic and security information about an IP address in a
quick and simple way. It uses various free public services to do so.

```
$ checkip 91.228.166.47
abuseipdb.com   domain: eset.sk, usage type: Data Center/Web Hosting/Transit
dns             name: skh1-webredir01-v.eset.com.
iptoasn.com     AS description: ESET-AS
maxmind.com     city: Bratislava, country: Slovakia (SK)
ping            0% packet loss, sent 5, recv 5, avg round-trip 9 ms
shodan.io       OS: n/a, 2 open ports: tcp/80 (nginx), tcp/443 (nginx)
urlscan.io      0 related URLs
virustotal.com  network: 91.228.164.0/22, SAN: www.eset.com, eset.com
Malicious       0% ✅
$ checkip 209.141.33.65
checkip: CheckDNS: lookup 209.141.33.65: nodename nor servname provided, or not known
abuseipdb.com   domain: buyvm.net, usage type: Data Center/Web Hosting/Transit
iptoasn.com     AS description: PONYNET - FranTech Solutions
maxmind.com     city: Las Vegas, country: United States (US)
ping            0% packet loss, sent 5, recv 5, avg round-trip 159 ms
shodan.io       OS: Debian, 2 open ports: tcp/22 (OpenSSH, 8.4p1 Debian 5), tcp/8443
urlscan.io      0 related URLs
virustotal.com  network: 209.141.32.0/19, SAN: 2020201.tk
Malicious       25% 🤏
$ checkip 218.92.0.158
checkip: CheckDNS: lookup 218.92.0.158: nodename nor servname provided, or not known
checkip: CheckShodan: GET https://api.shodan.io/shodan/host/218.92.0.158?key=REDACTED: 404 Not Found
abuseipdb.com   domain: chinatelecom.com.cn, usage type: n/a
iptoasn.com     AS description: CHINANET-BACKBONE No.31,Jin-rong Street
maxmind.com     city: Shanghai, country: China (CN)
ping            100% packet loss, sent 5, recv 0, avg round-trip 0 ms
urlscan.io      0 related URLs
virustotal.com  network: 218.92.0.0/16, SAN: n/a
Malicious       50% 🚫
```

The CLI tool also supports JSON output

```
# select only Info (1) and InfoSec (2) check types
checkip -j 118.25.6.39 | jq -r '.checks[] | select(.type == 1 or .type == 2) | "\(.malicious)\t\(.name)"' | sort
```

## Installation and configuration

To install the CLI tool

```
git clone git@github.com:jreisinger/checkip.git
cd checkip
make install
```

or download a [release](https://github.com/jreisinger/checkip/releases)
binary for your system and architecture.

For some checks to work you need to register and get an API (LICENSE) key.
Then create a `$HOME/.checkip.yaml` using your editor of choice

```
ABUSEIPDB_API_KEY: aaaaaaaabbbbbbbbccccccccddddddddeeeeeeeeffffffff11111111222222223333333344444444
MAXMIND_LICENSE_KEY: abcdef1234567890
SHODAN_API_KEY: aaaabbbbccccddddeeeeffff11112222
URLSCAN_API_KEY: abcd1234-a123-4567-678z-a2b3c4b5d6e7
VIRUSTOTAL_API_KEY: aaaaaaaabbbbbbbbccccccccddddddddeeeeeeeeffffffff1111111122222222
```

You can also use environment variables with the same names.

## Development

Checkip is easy to extend. If you want to add a new way to check an IP address,
just write a function of type
[Check](https://pkg.go.dev/github.com/jreisinger/checkip/check#Check).
