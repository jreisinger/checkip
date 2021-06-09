package checkip

import (
	"fmt"
	"net"
	"strings"

	"github.com/oschwald/geoip2-golang"
)

// Geo holds geographic location of an IP address from maxmind.com GeoIP database.
type Geo struct {
	Location []string
}

// Check fills in the geolocation data.
func (g *Geo) Check(ip net.IP) (bool, error) {
	licenseKey, err := GetConfigValue("GEOIP_LICENSE_KEY")
	if err != nil {
		return false, fmt.Errorf("can't download DB: %w", err)
	}

	file := "/var/tmp/GeoLite2-City.mmdb"
	url := "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=" + licenseKey + "&suffix=tar.gz"

	if err := Update(file, url, "tgz"); err != nil {
		return false, fmt.Errorf("can't update DB file: %v", err)
	}

	db, err := geoip2.Open(file)
	if err != nil {
		return false, fmt.Errorf("can't load DB file: %v", err)
	}
	defer db.Close()

	record, err := db.City(ip)
	if err != nil {
		return false, err
	}

	city := record.City.Names["en"]
	country := record.Country.Names["en"]
	isoCode := record.Country.IsoCode

	if city == "" {
		city = "city unknown"
	}
	if country == "" {
		country = "country unknown"
	}
	if isoCode == "" {
		isoCode = "ISO code unknown"
	}

	g.Location = append(g.Location, city)
	g.Location = append(g.Location, country)
	g.Location = append(g.Location, isoCode)

	return true, nil
}

// String returns the result of the check.
func (g *Geo) String() string {
	return fmt.Sprintf("%s", strings.Join(g.Location, ", "))
}
