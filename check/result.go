package check

import (
	"fmt"
	"log"
	"sort"

	"github.com/logrusorgru/aurora"
)

type Result struct {
	CheckName         string
	CheckType         Type
	IsIPaddrMalicious bool
	Data              Data
	Error             *ResultError
}

type Results []Result

func (rs Results) SortByName() {
	sort.Slice(rs, func(i, j int) bool {
		return rs[i].CheckName < rs[j].CheckName
	})
}

// Print prints condensed results to stdout.
func (rs Results) Print() error {

	var malicious, totalSec float64
	for _, r := range rs {
		if r.Error != nil {
			log.Print(r.Error.Error())
		}
		if r.CheckType == "Info" || r.CheckType == "InfoSec" {
			fmt.Printf("%-15s %s\n", r.CheckName, r.Data.String())
		}
		if r.CheckType == "Sec" || r.CheckType == "InfoSec" {
			totalSec++
			if r.IsIPaddrMalicious {
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
