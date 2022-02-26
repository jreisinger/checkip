# checkip
[![PkgGoDev](https://pkg.go.dev/badge/github.com/jreisinger/checkip)](https://pkg.go.dev/github.com/jreisinger/checkip)

Checkip is a CLI tool and library that provides generic and security information 
about an IP address in a quick way. It uses various free public services to do so.

```
$ checkip 91.228.166.47
abuseipdb.com  --> domain: eset.sk, usage type: Data Center/Web Hosting/Transit
dns mx         --> eset.sk: a.mx.eset.com, b.mx.eset.com
dns name       --> skh1-webredir01-v.eset.com
iptoasn.com    --> AS description: ESET-AS
maxmind.com    --> city: Bratislava, country: Slovakia (SK)
shodan.io      --> OS: n/a, 2 open ports: tcp/80 (nginx), tcp/443 (nginx)
urlscan.io     --> 0 related URLs
virustotal.com --> network: 91.228.164.0/22, SAN: www.eset.com, eset.com
Malicious      --> 0% ✅

$ checkip 209.141.33.65 2> /dev/null
abuseipdb.com  --> domain: buyvm.net, usage type: Data Center/Web Hosting/Transit
dns mx         --> buyvm.net: mail.frantech.ca
iptoasn.com    --> AS description: PONYNET - FranTech Solutions
maxmind.com    --> city: Las Vegas, country: United States (US)
shodan.io      --> OS: Debian, 1 open port: tcp/22 (OpenSSH, 8.4p1 Debian 5)
urlscan.io     --> 0 related URLs
virustotal.com --> network: 209.141.32.0/19, SAN: 2020201.tk
Malicious      --> 25% 🤏

$ checkip 218.92.0.158 2> /dev/null
abuseipdb.com  --> domain: chinatelecom.com.cn, usage type: n/a
dns mx         --> chinatelecom.com.cn: testmail.chinatelecom.com.cn
iptoasn.com    --> AS description: CHINANET-BACKBONE No.31,Jin-rong Street
maxmind.com    --> city: n/a, country: China (CN)
urlscan.io     --> 0 related URLs
virustotal.com --> network: 218.92.0.0/16, SAN: n/a
Malicious      --> 50% 🚫
```

Optionally it can also interact with the target IP address. NOTE: you should run
active checks (-a) only against your hosts or hosts you have
[permission](http://scanme.nmap.org/) to scan.

```
$ checkip -a 45.33.32.156 # scanme.nmap.org
Open TCP ports --> 22 (ssh), 80 (http), 9929 (nping-echo), 31337 (Elite)
Ping           --> 0% packet loss, sent 5, recv 5, avg round-trip 192 ms
abuseipdb.com  --> domain: linode.com, usage type: Data Center/Web Hosting/Transit
dns mx         --> linode.com: inbound-mail1.linode.com, inbound-mail3.linode.com
dns name       --> scanme.nmap.org
iptoasn.com    --> AS description: LINODE-AP Linode, LLC
maxmind.com    --> city: Fremont, country: United States (US)
shodan.io      --> OS: n/a, 3 open ports: tcp/22 (OpenSSH, 6.6.1p1 Ubuntu-2ubuntu2.13), tcp/80 (Apache httpd, 2.4.7), udp/123
urlscan.io     --> 0 related URLs
virustotal.com --> network: 45.33.0.0/17, SAN: n/a
Malicious      --> 0% ✅
```

The CLI tool also supports JSON output:

```
$ checkip -j 218.92.0.158 2> /dev/null | \
# select only Sec (1) or InfoSec (2) check types and show if considered malicious
jq -r '.checks[] | select(.type == 1 or .type == 2) | "\(.malicious)\t\(.name)"'
false	abuseipdb.com
true	blocklist.de
false	cinsscore.com
true	github.com/stamparm/ipsum
true	otx.alienvault.com
false	threatcrowd.org
false	urlscan.io
true	virustotal.com
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
function to `checks.Passive` or `checks.Active`
[variable](https://pkg.go.dev/github.com/jreisinger/checkip/checks#pkg-variables).
Then:

```
make run

git tag | sort -V
git tag -a v0.16.2 -m "improve docs"
git push -u origin --tags
```
