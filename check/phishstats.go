package check

import (
	"encoding/csv"
	"encoding/json"
	"net"
	"os"
	"strconv"
)

type phishstats struct {
	Score float64 // 0-2 likely, 2-4 suspicious, 4-6 phishing, 6-10 omg phishing!
	Url   string
}

func (ps phishstats) Summary() string {
	return ps.Url
}

func (ps phishstats) Json() ([]byte, error) {
	return json.Marshal(ps)
}

// PhishStats checks whether the ipaddr is involved in phishing according to
// https://phishstats.info/phish_score.csv.
func PhishStats(ipaddr net.IP) (Check, error) {
	result := Check{
		Description: "phishstats.info",
		Type:        InfoAndIsMalicious,
	}

	file, err := getCachePath("phish_score.csv")
	if err != nil {
		return result, err
	}
	url := "https://phishstats.info/phish_score.csv"
	if err := updateFile(file, url, ""); err != nil {
		return result, err
	}

	ps, err := getPhishStats(file, ipaddr)
	if err != nil {
		return result, err
	}
	result.IpAddrInfo = ps
	if ps.Score > 2 {
		result.IpAddrIsMalicious = true
	}

	return result, nil
}

func getPhishStats(csvFile string, ipaddr net.IP) (phishstats, error) {
	var ps phishstats

	f, err := os.Open(csvFile)
	if err != nil {
		return ps, err
	}

	csvReader := csv.NewReader(f)
	csvReader.Comment = '#'
	records, err := csvReader.ReadAll()
	if err != nil {
		return ps, err
	}

	for _, fields := range records {
		if ipaddr.String() == fields[3] {
			score, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				return ps, err
			}
			ps.Score = score
			ps.Url = fields[2]
			break
		}
	}

	return ps, nil
}
