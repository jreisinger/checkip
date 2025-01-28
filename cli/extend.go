// Package cli contains functions for running checks from command-line.
package cli

import (
	"fmt"
	"strings"

	"github.com/jreisinger/checkip/check"
)

// ExtPrintSummary add IpAddrInfo.Summary for IOCLoc check
func (rs Checks) ExtPrintSummary() string {
	res := ""
	cotation := ""
	for _, r := range rs {
		// To avoid "invalid memory address or nil pointer dereference"
		// runtime error and printing empty summary info.
		if r.IpAddrInfo == nil || r.IpAddrInfo.Summary() == "" {
			continue
		}

        desc := strings.ToLower(r.IpAddrInfo.Summary())
		switch {
		case cotation == "" && strings.Contains(desc, "data center"):
			cotation = "A"  // server
		case cotation == "" && strings.Contains(r.IpAddrInfo.Summary(), "open:"):
			cotation = "A"  // server
		case strings.Contains(desc, "vpn") || strings.Contains(desc, "avast"):
			cotation = "B"  // vpn, proxy
		case strings.Contains(desc, "mikrotik") || strings.Contains(desc, "fixed line"):
			cotation = "C"  // botnet
		case strings.Contains(desc, "mobile"):
			cotation = "D"  // mobile
		case strings.Contains(desc, "akamai") || strings.Contains(desc, "amazon"):
			cotation = "E"  // cdn
		}

		if r.Type == check.Info || r.Type == check.InfoAndIsMalicious {
			fmt.Printf("%-15s %s\n", r.Description, r.IpAddrInfo.Summary())
			if r.Description == "IOCLoc" {
				res = r.IpAddrInfo.Summary()
			}

		}

	}
	return fmt.Sprintf("%s [%s1]",res,cotation)
}
