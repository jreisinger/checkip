# checkip

Checkip is a CLI tool and library that provides generic and security information about an IP address in an easy way. It uses various [checks](https://pkg.go.dev/github.com/jreisinger/checkip/check) to do so.

```
$ checkip 218.92.0.158
abuseipdb.com   --> domain: chinatelecom.com.cn, usage type: Data Center/Web Hosting/Transit
db-ip.com       --> country: China (CN), city: Nanjing (Jiangning Qu), EU member: false
dns mx          --> chinatelecom.com.cn: testmail.chinatelecom.com.cn
dns name        --> n/a
iptoasn.com     --> AS description: CHINANET-BACKBONE No.31,Jin-rong Street
maxmind.com     --> country: China (CN), city: Caolin, EU member: false
phishstats.info --> n/a
shodan.io       --> OS: n/a, 2 open ports: tcp/22 (OpenSSH, 7.4), tcp/53
urlscan.io      --> 0 related URLs
virustotal.com  --> network: 218.92.0.0/16, SAN: n/a
Malicious       --> 50% (5/10) 🚫
```

You can get output also in JSON (`-j`). Here we select Sec (1) and InfoSec (2) check [type](https://pkg.go.dev/github.com/jreisinger/checkip#Type) and show if the check considers the IP address to be malicious.

```
$ checkip -j 218.92.0.158 | jq -r \
'.checks[] | select(.type==1 or .type==2) | "\(.malicious)\t\(.name)"'
true	abuseipdb.com
true	blocklist.de
false	cinsscore.com
false	firehol.org
true	github.com/stamparm/ipsum
true	otx.alienvault.com
false	phishstats.info
false	threatcrowd.org
false	urlscan.io
true	virustotal.com
```

Active checks (`-a`) directly interact with the IP address. And scan checks (`-s`) scan the IP address! You should run them only against your hosts or hosts you have permission to scan.

```
$ checkip -a -s 45.33.32.156 # scanme.nmap.org
checkip: Tls: dial tcp 45.33.32.156:443: connect: connection refused
Open TCP ports  --> 22 (ssh), 80 (http), 9929 (nping-echo), 31337 (Elite)
Ping            --> 0% packet loss, sent 5, recv 5, avg round-trip 169 ms
abuseipdb.com   --> domain: linode.com, usage type: Data Center/Web Hosting/Transit
db-ip.com       --> country: United States (US), city: Fremont, EU member: false
dns mx          --> linode.com: inbound-mail1.linode.com, inbound-mail3.linode.com
dns name        --> scanme.nmap.org
iptoasn.com     --> AS description: LINODE-AP Linode, LLC
maxmind.com     --> country: United States (US), city: Fremont, EU member: false
phishstats.info --> n/a
shodan.io       --> OS: n/a, 3 open ports: tcp/22 (OpenSSH, 6.6.1p1 Ubuntu-2ubuntu2.13), tcp/80 (Apache httpd, 2.4.7), udp/123
urlscan.io      --> 0 related URLs
virustotal.com  --> network: 45.33.0.0/17, SAN: n/a
Malicious       --> 0% (0/10) ✅
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

Checkip is easy to extend. If you want to add a new way to check an IP address, just write a function of type [Check](https://pkg.go.dev/github.com/jreisinger/checkip#Check). Add the function to `check.Passive`, `check.Active` or `check.Scan` [variable](https://pkg.go.dev/github.com/jreisinger/checkip/check#pkg-variables).

```
make run # test, install and run

git commit -m "improve tag docs" main.go

git tag | sort -V # or git ll
git tag -a v0.16.2 -m "improve docs"

git push --follow-tags
```
