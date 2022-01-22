# checkip

Checkip provides generic and security information about an IP address in a
quick way. It uses various free public services to do so.

```
$ checkip 91.228.166.47
abuseipdb.com  --> domain: eset.sk, usage type: Data Center/Web Hosting/Transit
dns mx         --> eset.sk: a.mx.eset.com, b.mx.eset.com
dns name       --> skh1-webredir01-v.eset.com
iptoasn.com    --> AS description: ESET-AS
maxmind.com    --> city: Bratislava, country: Slovakia (SK)
ping           --> 0% packet loss, sent 5, recv 5, avg round-trip 13 ms
shodan.io      --> OS: n/a, 2 open ports: tcp/80 (nginx), tcp/443 (nginx)
urlscan.io     --> 0 related URLs
virustotal.com --> network: 91.228.164.0/22, SAN: www.eset.com, eset.com
Malicious      --> 0% ✅
$ checkip 209.141.33.65 2> /dev/null
abuseipdb.com  --> domain: buyvm.net, usage type: Data Center/Web Hosting/Transit
dns mx         --> buyvm.net: mail.frantech.ca
iptoasn.com    --> AS description: PONYNET - FranTech Solutions
maxmind.com    --> city: Las Vegas, country: United States (US)
ping           --> 0% packet loss, sent 5, recv 5, avg round-trip 172 ms
shodan.io      --> OS: Debian, 1 open port: tcp/22 (OpenSSH, 8.4p1 Debian 5)
urlscan.io     --> 0 related URLs
virustotal.com --> network: 209.141.32.0/19, SAN: 2020201.tk
Malicious      --> 25% 🤏
$ checkip 218.92.0.158 2> /dev/null
abuseipdb.com  --> domain: chinatelecom.com.cn, usage type: n/a
dns mx         --> chinatelecom.com.cn: testmail.chinatelecom.com.cn
iptoasn.com    --> AS description: CHINANET-BACKBONE No.31,Jin-rong Street
maxmind.com    --> city: Lianyungang, country: China (CN)
ping           --> 100% packet loss, sent 5, recv 0, avg round-trip 0 ms
urlscan.io     --> 0 related URLs
virustotal.com --> network: 218.92.0.0/16, SAN: n/a
Malicious      --> 50% 🚫
```

The CLI tool also supports JSON output:

```
checkip -j 118.25.6.39 | \
# select only Sec (1) or InfoSec (2) check types
jq -r '.checks[] | select(.type == 1 or .type == 2) | "\(.malicious)\t\(.name)"'
```

## Installation and configuration

To install the CLI tool

```
git clone git@github.com:jreisinger/checkip.git
cd checkip
make install
```

or download a [release](https://github.com/jreisinger/checkip/releases)
binary (from under "Assets") for your system and architecture.

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
[Check](https://pkg.go.dev/github.com/jreisinger/checkip/check#Check). Add the
function to `checks.Default`
[variable](https://pkg.go.dev/github.com/jreisinger/checkip/checks#pkg-variables).
Then:

```
make run

git tag -a v0.16.2 -m "improve docs"
git push -u origin --tags
```
