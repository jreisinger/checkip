# checkip

CLI tool that finds out information about an IP address. Currently two types of information are supported:

* Geographic location using [GeoLite2 City database](https://dev.maxmind.com/geoip/geoip2/geolite2/). It takes care of downloading the database if it's not present or it's outdated. You need to set the `GEOIP_LICENSE_KEY` environment variable so it can download the database. Read [this](https://dev.maxmind.com/geoip/geoip2/geolite2/#Download_Access) for how to get the license key (it's free).
* ASN using [https://iptoasn.com/](iptoasn) API.
* DNS name using [net.LookupAddr](https://golang.org/pkg/net/#LookupAddr).

Installation:

```
make install
```

Usage:

```
> checkip 1.1.1.1
Geo (maxmind.com)   : city unknown, Australia, AU
ASN (iptoasn.com)   : 7497, 1.1.1.1 - 1.1.1.1, CSTNET-AS-AP Computer Network Information Center, CN
DNS (net.LookupAddr): one.one.one.one.

> checkip $(dig +short google.com)
Geo (maxmind.com)   : Mountain View, United States, US
ASN (iptoasn.com)   : 15169, 216.58.192.0 - 216.58.223.255, GOOGLE - Google LLC, US
DNS (net.LookupAddr): prg03s01-in-f14.1e100.net., prg03s01-in-f78.1e100.net.

> checkip $(curl -s util.reisinge.net/addr) # your own IP address
Geo (maxmind.com)   : Bratislava, Slovakia, SK
ASN (iptoasn.com)   : 15962, 109.230.0.0 - 109.230.63.255, OSK-DNI Slovakia, SK
DNS (net.LookupAddr): 109-230-20-210.dynamic.orange.sk.
```
