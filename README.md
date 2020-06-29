[![Build Status](https://travis-ci.org/jreisinger/checkip.svg?branch=master)](https://travis-ci.org/jreisinger/checkip)

# checkip

CLI tool that finds out information about an IP address. Currently these types of information are provided:

* ASN using [iptoasn](https://iptoasn.com/) API.
* DNS name using [net.LookupAddr](https://golang.org/pkg/net/#LookupAddr).
* Geographic location using [GeoLite2 City database](https://dev.maxmind.com/geoip/geoip2/geolite2/). Read [this](https://dev.maxmind.com/geoip/geoip2/geolite2/#Download_Access) for how to get the license key (it's free).

## Installation

Download the latest [release](https://github.com/jreisinger/checkip/releases) for your operating system and architecture.

## Usage examples

```
$ checkip 255.255.255.256 # you must supply a valid IP address
invalid IP address: 255.255.255.256

$ checkip 1.1.1.1
Geo: city unknown, Australia, AU
DNS: one.one.one.one.
ASN: 7497, 1.1.1.1 - 1.1.1.1, CSTNET-AS-AP Computer Network Information Center, CN
```

## Development

```
vim main.go
make install # version defaults to "dev"

make release # you'll find releases in releases/ directory
```

Builds are done inside Docker container. Once you push to GitHub Travis will
try and build a release for you and publish it on GitHub.