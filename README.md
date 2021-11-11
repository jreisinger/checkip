# checkip

Checkip checks an IP address using various public services. It provides
generic and security information in a quick and simple way. It can be run as
a CLI tool or as JSON API server.

<img src="checkip.png" width="800">

## Installation and configuration

To install the CLI/API tool

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
