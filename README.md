# checkip

Checkip is a CLI tool and library that provides generic and security information
about an IP address in quick way. It uses various free public services to do so.

```
$ checkip 218.92.0.158
abuseipdb.com  --> domain: chinatelecom.com.cn, usage type: Data Center/Web Hosting/Transit
db-ip.com      --> country: China (CN), city: Nanjing (Jiangning Qu), EU member: false
dns mx         --> chinatelecom.com.cn: testmail.chinatelecom.com.cn
dns name       --> n/a
iptoasn.com    --> AS description: CHINANET-BACKBONE No.31,Jin-rong Street
maxmind.com    --> country: China (CN), city: Caolin, EU member: false
shodan.io      --> OS: n/a, 2 open ports: tcp/22 (OpenSSH, 7.4), tcp/53
urlscan.io     --> 0 related URLs
virustotal.com --> network: 218.92.0.0/16, SAN: n/a
Malicious      --> 63% (5/8) 🚫
```

The CLI tool also supports JSON output.

```
$ checkip -j 218.92.0.158 | \
# Select Sec (1) and InfoSec (2) check types that returned some info
# (i.e. they worked) and show if they consider the IP address malicious.
jq -r '.checks[] | select((.type == 1 or .type == 2) and .info != null) | "\(.malicious)\t\(.name)"'
true	abuseipdb.com
true	blocklist.de
false	cinsscore.com
true	github.com/stamparm/ipsum
true	otx.alienvault.com
false	threatcrowd.org
false	urlscan.io
true	virustotal.com
```

Active checks (`-a`) interact with the IP address. You should only run
them against your hosts or hosts you have permission to scan.

```
$ checkip -a 45.33.32.156 # scanme.nmap.org
Open TCP ports --> 22 (ssh), 80 (http), 9929 (nping-echo), 31337 (Elite)
Ping           --> 0% packet loss, sent 5, recv 5, avg round-trip 168 ms
abuseipdb.com  --> domain: linode.com, usage type: Data Center/Web Hosting/Transit
db-ip.com      --> country: United States (US), city: Fremont, EU member: false
dns mx         --> linode.com: inbound-mail1.linode.com, inbound-mail3.linode.com
dns name       --> scanme.nmap.org
iptoasn.com    --> AS description: LINODE-AP Linode, LLC
maxmind.com    --> country: United States (US), city: Fremont, EU member: false
shodan.io      --> OS: n/a, 3 open ports: tcp/22 (OpenSSH, 6.6.1p1 Ubuntu-2ubuntu2.13), tcp/80 (Apache httpd, 2.4.7), udp/123
urlscan.io     --> 0 related URLs
virustotal.com --> network: 45.33.0.0/17, SAN: n/a
Malicious      --> 0% (0/8) ✅
```

## Installation

To install the CLI tool

```
# optional; to install inside a container
docker run --rm -it golang /bin/bash

go install github.com/jreisinger/checkip@latest
```

or download a [release](https://github.com/jreisinger/checkip/releases)
binary (from under "Assets") for your system and architecture.

## Configuration

For some checks to work you need to register and get an API (LICENSE) key. See
the service web site for how to do that.

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

Checkip is easy to extend. If you want to add a new way to check an IP address,
just write a function of type [Check][1]. Add the function to `checks.Passive`
or `checks.Active` [variable][2].

```
make run # test, install and run

git commit -m "improve tag docs" main.go

git tag | sort -V
git tag -a v0.16.2 -m "improve docs"

git push --follow-tags
```

[1]: https://pkg.go.dev/github.com/jreisinger/checkip/check#Check
[2]: https://pkg.go.dev/github.com/jreisinger/checkip/checks#pkg-variables
