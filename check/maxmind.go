package check

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/oschwald/geoip2-golang"
)

type maxmind struct {
	City    string `json:"city"`
	Country string `json:"country"`
	IsoCode string `json:"iso_code"`
	IsInEU  bool   `json:"is_in_eu"`
}

func (m maxmind) Summary() string {
	// Get just non-empty strings.
	var parts []string
	for _, s := range []string{m.City, m.Country} {
		if strings.TrimSpace(s) != "" {
			parts = append(parts, s)
		}
	}

	return strings.Join(parts, ", ")
}

func (m maxmind) Json() ([]byte, error) {
	return json.Marshal(m)
}

// MaxMind gets geolocation data from maxmind.com's GeoLite2-City.mmdb.
func MaxMind(ip net.IP) (Check, error) {
	result := Check{
		Description: "maxmind.com",
		Type:        Info,
	}

	licenseKey, err := getConfigValue("MAXMIND_LICENSE_KEY")
	if err != nil {
		return result, newCheckError(err)
	}
	if licenseKey == "" {
		result.MissingCredentials = "MAXMIND_LICENSE_KEY"
		return result, nil
	}

	// file := "/var/tmp/GeoLite2-City.mmdb"
	file, err := getCachePath("GeoLite2-City.mmdb")
	if err != nil {
		return result, err
	}

	url := "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=" + licenseKey + "&suffix=tar.gz"

	if err := updateFile(file, url, "tgz"); err != nil {
		return result, newCheckError(err)
	}

	db, err := geoip2.Open(file)
	if err != nil {
		return result, newCheckError(fmt.Errorf("can't load DB file %s: %v", file, err))
	}
	defer db.Close()

	geo, err := db.City(ip)
	if err != nil {
		return result, newCheckError(err)
	}

	result.IpAddrInfo = maxmind{
		City:    geo.City.Names["en"],
		Country: geo.Country.Names["en"],
		IsoCode: geo.Country.IsoCode,
		IsInEU:  geo.Country.IsInEuropeanUnion,
	}

	return result, nil
}
