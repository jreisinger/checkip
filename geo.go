package checkip

import (
	"fmt"
	"net"

	"github.com/oschwald/geoip2-golang"
)

// Geo holds geographic location of an IP address from maxmind.com GeoIP database.
type Geo struct {
	City, Country, IsoCode string
}

func (g *Geo) Name() string { return "maxmind.com" }

// Check fills in the geolocation data. The data is taken from
// GeoLite2-City.mmdb file that gets downloaded and regularly updated.
func (g *Geo) Check(ip net.IP) error {
	licenseKey, err := getConfigValue("MAXMIND_LICENSE_KEY")
	if err != nil {
		return fmt.Errorf("can't download DB: %w", err)
	}

	file := "/var/tmp/GeoLite2-City.mmdb"
	url := "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=" + licenseKey + "&suffix=tar.gz"

	if err := updateFile(file, url, "tgz"); err != nil {
		return fmt.Errorf("can't update DB file: %v", err)
	}

	db, err := geoip2.Open(file)
	if err != nil {
		return fmt.Errorf("can't load DB file: %v", err)
	}
	defer db.Close()

	record, err := db.City(ip)
	if err != nil {
		return err
	}

	g.City = record.City.Names["en"]
	g.Country = record.Country.Names["en"]
	g.IsoCode = record.Country.IsoCode

	return nil
}

// Info returns interesting information from the check.
func (g *Geo) Info() string {
	if g.City == "" {
		g.City = "city unknown"
	}
	if g.Country == "" {
		g.Country = "country unknown"
	}
	if g.IsoCode == "" {
		g.IsoCode = "ISO code unknown"
	}
	return fmt.Sprintf("%s, %s (%s)", g.City, g.Country, g.IsoCode)
}
