# checkip

`checkip` is a CLI tool and library that checks an IP address using various
public services. It provides generic and security information in a simple and
quick way.

<img src="checkip.png" width="800">

The CLI tool can also print all data in JSON format so you can pick what you
want to see, e.g.:

```
checkip -j 218.92.0.158 | \
jq -r '.[] | select(.Type=="Sec" or .Type=="InfoSec") | "\(.IsMalicious)\t\(.Name)"' | sort
false	abuseipdb.com
false	cinsscore.com
false	threatcrowd.org
false	virustotal.com
true	blocklist.de
true	github.com/stamparm/ipsum
true	otx.alienvault.com
```

## Installation and configuration

To install the CLI tool

```
git clone git@github.com:jreisinger/checkip.git
cd checkip
make install
```

or download a [release](https://github.com/jreisinger/checkip/releases)
binary for your system and architecture.

For some checks to work you need to register and get an API (LICENSE) key.
Then create a `$HOME/.checkip.yaml` using your editor of choice

```
ABUSEIPDB_API_KEY: aaaaaaaabbbbbbbbccccccccddddddddeeeeeeeeffffffff11111111222222223333333344444444
MAXMIND_LICENSE_KEY: abcdef1234567890
SHODAN_API_KEY: aaaabbbbccccddddeeeeffff11112222
VIRUSTOTAL_API_KEY: aaaaaaaabbbbbbbbccccccccddddddddeeeeeeeeffffffff1111111122222222
```

You can also use environment variables with the same names.

## Development and releasing

An IP address is checked by running one or more
[Checker](https://pkg.go.dev/github.com/jreisinger/checkip#Checker)s. To add
a new way for checking IP addresses just implement the
[InfoChecker](https://pkg.go.dev/github.com/jreisinger/checkip#InfoChecker),
[SecChecker](https://pkg.go.dev/github.com/jreisinger/checkip#SecChecker) or
[InfoSecChecker](https://pkg.go.dev/github.com/jreisinger/checkip#InfoSecChecker)
interface and add it to `cmd/checkip.go`. Then

```
make run # see the picture above
```

If you are satisfied commit, push and add new [tag](https://reisinge.net/notes/prog/git#tags)

```
git tag -a v0.6.9 -m "goreleaser with GitHub Actions"
git push --tags
```

GitHub Actions with goreleaser will make and publish the release.