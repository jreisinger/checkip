# checkip

Find out geographic location of an IP address.

Installation:

```
> go install
```

Usage:

```
> checkip 1.1.1.1
, Australia, AU

> checkip $(dig +short reisinge.net)
Frankfurt am Main, Germany, DE

> checkip $(curl -s util.reisinge.net/addr)
Stupava, Slovakia, SK
```
