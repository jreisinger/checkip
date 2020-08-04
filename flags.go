package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func handleFlags() {

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s [flags] <ipaddr>\n", os.Args[0])
		flag.PrintDefaults()
	}

	version := flag.Bool("version", false, "version")

	flag.Parse()

	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}

	if len(flag.Args()) != 1 {
		log.Fatalf("missing IP address to check")
	}
}
