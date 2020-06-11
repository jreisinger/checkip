[![Build Status](https://travis-ci.org/jreisinger/checkip.svg?branch=master)](https://travis-ci.org/jreisinger/checkip)

# checkip

CLI tool that finds out information about an IP address. Currently these types of information are supported:

* Geographic location using [GeoLite2 City database](https://dev.maxmind.com/geoip/geoip2/geolite2/). It takes care of downloading the database if it's not present or it's outdated. You need to set the `GEOIP_LICENSE_KEY` environment variable so it can download the database. Read [this](https://dev.maxmind.com/geoip/geoip2/geolite2/#Download_Access) for how to get the license key (it's free).
* ASN using [iptoasn](https://iptoasn.com/) API.
* DNS name using [net.LookupAddr](https://golang.org/pkg/net/#LookupAddr).

## Installation

Download the latest [release](https://github.com/jreisinger/checkip/releases) for your operating system and architecture.

## Usage examples

```
> checkip 255.255.255.256 # you must supply a valid IP address
invalid IP address: 255.255.255.256

> checkip 1.1.1.1
Geo (maxmind.com)   : city unknown, Australia, AU
ASN (iptoasn.com)   : 7497, 1.1.1.1 - 1.1.1.1, CSTNET-AS-AP Computer Network Information Center, CN
DNS (net.LookupAddr): one.one.one.one.

> checkip $(dig +short google.com)
Geo (maxmind.com)   : Mountain View, United States, US
ASN (iptoasn.com)   : 15169, 216.58.192.0 - 216.58.223.255, GOOGLE - Google LLC, US
DNS (net.LookupAddr): prg03s01-in-f14.1e100.net., prg03s01-in-f78.1e100.net.

> checkip $(curl -s util.reisinge.net/addr) # this will check your own IP address
Geo (maxmind.com)   : Frankfurt am Main, Germany, DE
ASN (iptoasn.com)   : 16509, 52.56.0.0 - 52.60.255.255, AMAZON-02 - Amazon.com, Inc., US
DNS (net.LookupAddr): ec2-52-59-197-254.eu-central-1.compute.amazonaws.com
```

## Development

```
vim main.go
make install
make release # you'll find releases in releases/ directory
# if you push to GitHub Travis will build a release for you
```
