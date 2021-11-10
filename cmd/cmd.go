// Checkip quickly finds information about an IP address from a CLI.
package cmd

import (
	"fmt"
	"log"
	"sort"

	checkip "github.com/jreisinger/checkip/pkg"
	"github.com/logrusorgru/aurora"
)

type byName []checkip.Result

func (x byName) Len() int           { return len(x) }
func (x byName) Less(i, j int) bool { return x[i].Name < x[j].Name }
func (x byName) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// Print prints condensed results to stdout.
func Print(results []checkip.Result) error {
	sort.Sort(byName(results))

	var malicious, totalSec float64
	for _, r := range results {
		if r.Err != nil {
			log.Print(r.ErrMsg)
		}
		if r.Type == "Info" || r.Type == "InfoSec" {
			fmt.Printf("%-15s %s\n", r.Name, r.Info)
		}
		if r.Type == "Sec" || r.Type == "InfoSec" {
			totalSec++
			if r.IsMalicious {
				malicious++
			}
		}
	}
	probabilityMalicious := malicious / totalSec

	var msg string
	switch {
	case probabilityMalicious <= 0.15:
		msg = fmt.Sprint(aurora.Green("Malicious"))
	case probabilityMalicious <= 0.50:
		msg = fmt.Sprint(aurora.Yellow("Malicious"))
	default:
		msg = fmt.Sprint(aurora.Red("Malicious"))
	}

	_, err := fmt.Printf("%s\t%.0f%% (%d/%d)\n", msg, probabilityMalicious*100, int(malicious), int(totalSec))
	return err
}
