# checkip

`checkip` is a CLI tool and library that checks an IP address using various
public services.

![](checkip.png)

To install the CLI tool

```
git clone git@github.com:jreisinger/checkip.git
cd checkip
make install
```

or download a [release](https://github.com/jreisinger/checkip/releases)
binary for your system and architecture.

See `cmd/checkip.go` for how to use checkip as a library.

## Configuration

For some checks to work you need to register and get an API (LICENSE) key.
Then create a `$HOME/.checkip.yaml` using your editor of choice:

```
ABUSEIPDB_API_KEY: aaaaaaaabbbbbbbbccccccccddddddddeeeeeeeeffffffff11111111222222223333333344444444
GEOIP_LICENSE_KEY: abcdef1234567890
SHODAN_API_KEY: aaaabbbbccccddddeeeeffff11112222
VIRUSTOTAL_API_KEY: aaaaaaaabbbbbbbbccccccccddddddddeeeeeeeeffffffff1111111122222222
```

You can also use environment variables with the same names.

## Adding Checkers

To add a new service for checking IP addresses (i.e. a Checker) just
implement the
[checkip.Checker](https://pkg.go.dev/github.com/jreisinger/checkip#Checker)
interface and add it to `cmd/checkip.go`. Then:

```
make run
```

## Releasing

```
make release # build binaries in cmd
```

See [git tags](https://reisinge.net/notes/prog/git#tags).
