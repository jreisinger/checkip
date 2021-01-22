package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jreisinger/checkip/check"
)

// Version is the default version of checkip.
var Version = "dev"

func main() {
	log.SetPrefix(os.Args[0] + ": ") // prefix program name
	log.SetFlags(0)                  // no timestamp in error messages

	flags, err := ParseFlags()
	if err != nil {
		log.Fatalf("%v", err)
	}

	if flags.Version {
		fmt.Println(Version)
		os.Exit(0)
	}

	checks := flags.ChecksToRun

	// No -checks flag means run all available checks.
	if len(checks) == 0 {
		for _, chk := range check.GetAvailable() {
			checks = append(checks, chk)
		}
	}

	ch := make(chan string)
	for _, chk := range checks {
		go check.Run(chk, flags.IPaddr, ch)
	}
	for range checks {
		fmt.Print(<-ch)
	}
}
