[![Build Status](https://travis-ci.org/jreisinger/checkip.svg?branch=master)](https://travis-ci.org/jreisinger/checkip)

# checkip

CLI tool that finds out information about an IP address. Currently these types of information are provided:

* Confidence score that the IP address is malicious using [Abuse IP DB](https://www.abuseipdb.com). You need to [register](https://www.abuseipdb.com/register?plan=free) to get the API key (it's free).
* ASN using [iptoasn](https://iptoasn.com/) API.
* DNS name using [net.LookupAddr](https://golang.org/pkg/net/#LookupAddr).
* Geographic location using [GeoLite2 City database](https://dev.maxmind.com/geoip/geoip2/geolite2/). Read [this](https://dev.maxmind.com/geoip/geoip2/geolite2/#Download_Access) for how to get the license key (it's free).

## Installation

Download the latest [release](https://github.com/jreisinger/checkip/releases) for your operating system and architecture.

## Usage

```
$ checkip 255.255.255.256 # you must supply a valid IP address
invalid IP address: 255.255.255.256

$ checkip 1.1.1.1
GEO        city unknown, Australia, AU
DNS        one.one.one.one.
ASN        13335, 1.1.1.0 - 1.1.1.255, CLOUDFLARENET - Cloudflare, Inc., US
AbuseIPDB  malicious with 0% confidence (cloudflare.com)
```

## Development

```
vim main.go
make install # version defaults to "dev"

make release # you'll find releases in releases/ directory
```

Builds are done inside Docker container. Once you push to GitHub Travis will
try and build a release for you and publish it on GitHub.