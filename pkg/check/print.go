package check

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"log"
)

// Print prints condensed results to stdout.
func Print(results Results) error {

	var malicious, totalSec float64
	for _, r := range results {
		if r.Error != nil {
			log.Print(r.Error.Error())
		}
		if r.Type == "Info" || r.Type == "InfoSec" {
			fmt.Printf("%-15s %s\n", r.Name, r.Data.String())
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
