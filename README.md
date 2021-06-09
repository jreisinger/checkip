`checkip` is a CLI tool and library that checks an IP address using various
public services.

# Installation and usage

To install the CLI tool:

```
git clone git@github.com:jreisinger/checkip.git
cd checkip
make install
```

To use the CLI tool:

![](checkip.png)

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

## Development

To add a new service for checking IP addresses just implement the
checkip.Checker interface and then:

```
make run
```