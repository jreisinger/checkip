package check

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/logrusorgru/aurora"
)

type Result struct {
	Name            string
	Type            Type
	IPaddrMalicious bool
	Data            Data
	Error           *ResultError
}

type Results []Result

// PrintJSON prints all results in JSON.
func (rs Results) PrintJSON() {
	enc := json.NewEncoder(os.Stdout)
	if err := enc.Encode(rs); err != nil {
		log.Fatal(err)
	}
}

// SortByName sorts results by the checker name.
func (rs Results) SortByName() {
	sort.Slice(rs, func(i, j int) bool {
		return rs[i].Name < rs[j].Name
	})
}

// PrintInfo prints results from Info and InfoSec checkers.
func (rs Results) PrintInfo() {
	for _, r := range rs {
		if r.Error != nil {
			log.Print(r.Error.Error())
		}
		if r.Type == "Info" || r.Type == "InfoSec" {
			fmt.Printf("%-15s %s\n", r.Name, r.Data.String())
		}
	}
}

// PrintProbabilityMalicious prints the probability the IP address is malicious.
func (rs Results) PrintProbabilityMalicious() {
	var msg string
	switch {
	case rs.probabilityMalicious() <= 0.15:
		msg = fmt.Sprint(aurora.Green("Malicious"))
	case rs.probabilityMalicious() <= 0.50:
		msg = fmt.Sprint(aurora.Yellow("Malicious"))
	default:
		msg = fmt.Sprint(aurora.Red("Malicious"))
	}

	fmt.Printf("%s\t%.0f%%\n", msg, rs.probabilityMalicious()*100)
}

func (rs Results) probabilityMalicious() float64 {
	var malicious, totalSec float64
	for _, r := range rs {
		if r.Type == "Sec" || r.Type == "InfoSec" {
			totalSec++
			if r.IPaddrMalicious {
				malicious++
			}
		}
	}
	return malicious / totalSec
}
