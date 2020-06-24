[![Build Status](https://travis-ci.org/jreisinger/checkip.svg?branch=master)](https://travis-ci.org/jreisinger/checkip)

# checkip

CLI tool that finds out information about an IP address. Currently these types of information are provided:

* DNS name using [net.LookupAddr](https://golang.org/pkg/net/#LookupAddr).
* ASN using [iptoasn](https://iptoasn.com/) API.
* Geographic location using [GeoLite2 City database](https://dev.maxmind.com/geoip/geoip2/geolite2/). Read [this](https://dev.maxmind.com/geoip/geoip2/geolite2/#Download_Access) for how to get the license key (it's free).

## Installation

Download the latest [release](https://github.com/jreisinger/checkip/releases) for your operating system and architecture.

## Usage examples

```
$ checkip 255.255.255.256 # you must supply a valid IP address
invalid IP address: 255.255.255.256

$ checkip 1.1.1.1
DNS: one.one.one.one.
ASN: 7497, 1.1.1.1 - 1.1.1.1, CSTNET-AS-AP Computer Network Information Center, CN
Geo: city unknown, Australia, AU

$ checkip $(dig +short google.com)
DNS: bud02s26-in-f14.1e100.net.
ASN: 15169, 172.217.0.0 - 172.217.31.255, GOOGLE - Google LLC, US
Geo: city unknown, United States, US

$ checkip $(curl -s util.reisinge.net/addr) # this will check your own IP address
DNS: ec2-52-59-197-254.eu-central-1.compute.amazonaws.com.
ASN: 16509, 52.56.0.0 - 52.60.255.255, AMAZON-02 - Amazon.com, Inc., US
Geo: Frankfurt am Main, Germany, DE
```

## Development

```
vim main.go
make install
make release # you'll find releases in releases/ directory
# Travis will build a release for you and publish it on GitHub
```