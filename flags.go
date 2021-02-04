package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/jreisinger/checkip/check"
)

// Flags are all the available CLI flags (including arguments).
type Flags struct {
	Version     bool
	ChecksToRun checksToRun
	IPaddrs     []net.IP
	JSON        bool
}

// ParseFlags validates the flags and parses them into Flags.
func ParseFlags() (Flags, error) {
	var flags Flags

	f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	f.BoolVar(&flags.Version, "version", false, "print version")
	f.BoolVar(&flags.JSON, "json", false, "print output in JSON")
	f.Var(&flags.ChecksToRun, "check", "run only `CHECK[,CHECK,...]` instead of all checks")

	f.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s [FLAGS] [IPADDR[ IPADDR ...]]\n", os.Args[0])
		f.PrintDefaults()
	}

	err := f.Parse(os.Args[1:])
	if err != nil {
		return flags, err
	}

	if flags.Version {
		return flags, nil
	}

	// Read from STDIN.
	if len(f.Args()) == 0 {
		// return flags, fmt.Errorf("missing IP address to check")
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			addr := net.ParseIP(s.Text())
			if addr == nil {
				return flags, fmt.Errorf("invalid IP address: %v", s.Text())
			}
			flags.IPaddrs = append(flags.IPaddrs, addr)
		}
	}

	// Read CLI arguments.
	for _, arg := range f.Args() {
		addr := net.ParseIP(arg)
		if addr == nil {
			return flags, fmt.Errorf("invalid IP address: %v", arg)
		}
		flags.IPaddrs = append(flags.IPaddrs, addr)
	}

	return flags, err
}

// checksToRun can be used multiple times and can take multiple values separated
// by a comma. It contains the checks to run selected via -check.
type checksToRun []check.Check

func (a *checksToRun) String() string {
	return fmt.Sprintf("%s", *a)
}

func (a *checksToRun) Set(value string) error {
	requestedCheckNames := strings.Split(value, ",")
	for _, reqChkName := range requestedCheckNames {
		chk, ok := isAvailable(reqChkName)
		if !ok {
			log.Fatalf("unknown check: %s\n", reqChkName)
		}
		*a = append(*a, chk)
	}
	return nil
}

func isAvailable(checkName string) (check.Check, bool) {
	if checkName == "" {
		return nil, false
	}

	checkName = strings.TrimSpace(checkName)
	checkName = strings.ToLower(checkName)

	for _, chk := range check.GetAvailable() {
		if strings.HasPrefix(strings.ToLower(chk.Name()), checkName) {
			return chk, true
		}
	}

	return nil, false
}
