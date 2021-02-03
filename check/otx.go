package check

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/jreisinger/checkip/util"
)

// OTX holds IP address reputation data from otx.alienvault.com.
type OTX struct {
	Reputation struct {
		ThreatScore int         `json:"threat_score"`
		Counts      interface{} `json:"counts"`
		FirstSeen   string      `json:"first_seen"`
		LastSeen    string      `json:"last_seen"`
	} `json:"reputation"`
}

// Do gets data from https://otx.alienvault.com/api. It returns false when
// threat score is higher than two.
func (otx *OTX) Do(ipaddr net.IP) (bool, error) {
	otxurl := fmt.Sprintf("https://otx.alienvault.com/api/v1/indicators/IPv4/%s/reputation", ipaddr.String())
	baseURL, err := url.Parse(otxurl)
	if err != nil {
		return false, err
	}

	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		return false, err
	}

	client := NewHTTPClient(5 * time.Second)
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("%s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(otx); err != nil {
		return false, err
	}

	if otx.isNotOK() {
		return false, nil
	}

	return true, nil
}

func (otx *OTX) isNotOK() bool {
	return otx.Reputation.ThreatScore > 2
}

// Name returns the name of the check.
func (otx *OTX) Name() string {
	return fmt.Sprint("OTX")
}

// String returns the result of the check.
func (otx *OTX) String() string {
	var activities []string

	if otx.Reputation.Counts != nil {
		counts := otx.Reputation.Counts.(map[string]interface{})
		for activity, n := range counts {
			activities = append(activities, activity+" - "+fmt.Sprint(n))
		}
	}

	score := strconv.Itoa(otx.Reputation.ThreatScore)
	if otx.isNotOK() {
		score = util.Highlight(score)
	}

	return fmt.Sprintf("threat score %s (seen %s - %s)",
		score,
		parseTime(otx.Reputation.FirstSeen),
		parseTime(otx.Reputation.LastSeen),
	)
}

func parseTime(value string) string {
	if value == "" {
		return "date unknown"
	}
	inlayout := "2006-01-02T15:04:05"
	outlayout := "2006-01-02"
	t, err := time.Parse(inlayout, value)
	if err != nil {
		log.Fatal(err)
	}
	return t.Format(outlayout)
}
