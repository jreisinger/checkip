package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/jreisinger/checkip"
	"github.com/logrusorgru/aurora"
)

var s = flag.Bool("s", false, "use only checkers that tell whether ipaddr is suspicious")

func main() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Printf("Usage: %s <ipaddr>\n", os.Args[0])
		os.Exit(0)
	}

	ipaddr := net.ParseIP(flag.Arg(0))
	if ipaddr == nil {
		fmt.Fprintf(os.Stderr, "%s: wrong IP address: %s\n", os.Args[0], flag.Arg(0))
		os.Exit(1)
	}

	// secCheckers can tell you wether the IP address is suspicious.
	secCheckers := map[string]checkip.Checker{
		"abuseipdb.com":             &checkip.AbuseIPDB{},
		"otx.alienvault.com":        &checkip.OTX{},
		"github.com/stamparm/ipsum": &checkip.IPsum{},
		"shodan.io":                 &checkip.Shodan{},
		"threatcrowd.org":           &checkip.ThreatCrowd{},
		"virustotal.com":            &checkip.VirusTotal{},
	}

	// infoCheckers just give you information about an IP address.
	infoCheckers := map[string]checkip.Checker{
		"iptoasn.com":          &checkip.AS{},
		"net.LookupAddr":       &checkip.DNS{},
		"maxmind.com GeoLite2": &checkip.Geo{},
	}

	if !*s {
		for k, v := range infoCheckers {
			secCheckers[k] = v
		}
	}

	ch := make(chan string)
	format := "%-25s %s"
	for name, checker := range secCheckers {
		go func(checker checkip.Checker, name string) {
			ok, err := checker.Check(ipaddr)
			switch {
			case err != nil:
				ch <- fmt.Sprintf(format, name, aurora.Gray(11, err.Error()))
			case !ok:
				ch <- fmt.Sprintf(format, name, aurora.Magenta(checker.String()))
			default:
				ch <- fmt.Sprintf(format, name, checker)
			}
		}(checker, name)
	}
	for range secCheckers {
		fmt.Println(<-ch)
	}
}
