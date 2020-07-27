package virustotal

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
)

// New creates VirusTotal.
func New() *VirusTotal {
	return &VirusTotal{}
}

// ForIP fills in data for a given IP address.
func (t *VirusTotal) ForIP(ipaddr net.IP) error {
	apiKey := os.Getenv("VIRUSTOTAL_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("can't call API: environment variable VIRUSTOTAL_API_KEY is not set")
	}

	// curl --header "x-apikey:$VIRUSTOTAL_API_KEY" https://www.virustotal.com/api/v3/ip_addresses/1.1.1.1

	baseURL, err := url.Parse("https://www.virustotal.com/api/v3/ip_addresses/" + ipaddr.String())
	if err != nil {
		return err
	}

	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		return err
	}

	// Set request headers.
	req.Header.Set("x-apikey", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(t); err != nil {
		return err
	}

	return nil
}
