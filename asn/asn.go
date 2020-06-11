package asn

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

// Autonomous System
type AS struct {
	CountryCode string `json:"as_country_code"`
	Number      int    `json:"as_number"`
	Description string `json:"as_description"`
	FirsIP      net.IP `json:"first_ip"`
	LastIP      net.IP `json:"last_ip"`
}

func New() *AS {
	return &AS{}
}

func (a *AS) ForIP(ipaddr net.IP) error {
	resp, err := http.Get("https://api.iptoasn.com/v1/as/ip/" + ipaddr.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("search asn failed: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(a); err != nil {
		return err
	}

	return nil
}
