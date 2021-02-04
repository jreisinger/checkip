package check

import (
	"fmt"
	"net"

	"github.com/jreisinger/checkip/util"
	"github.com/oschwald/geoip2-golang"
)

// Geo holds geographic position of an IP address from MaxMind's GeoIP database.
type Geo struct {
	City    string
	Country string
	ISOCode string
}

// Do fills in the geolocation data.
func (g *Geo) Do(ip net.IP) (bool, error) {
	licenseKey, err := util.GetConfigValue("GEOIP_LICENSE_KEY")
	if err != nil {
		return false, fmt.Errorf("can't download DB: %w", err)
	}

	file := "/var/tmp/GeoLite2-City.mmdb"
	url := "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=" + licenseKey + "&suffix=tar.gz"

	if err := util.Update(file, url, "tgz"); err != nil {
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

	g.City = record.City.Names["en"]
	g.Country = record.Country.Names["en"]
	g.ISOCode = record.Country.IsoCode

	if g.City == "" {
		g.City = "city unknown"
	}
	if g.Country == "" {
		g.Country = "country unknown"
	}
	if g.ISOCode == "" {
		g.ISOCode = "ISO code unknown"
	}

	return true, nil
}

// Name returns the name of the check.
func (g *Geo) Name() string {
	return fmt.Sprint("Geo")
}

// String returns the result of the check.
func (g *Geo) String() string {
	return fmt.Sprintf("%s, %s, %s", g.City, g.Country, g.ISOCode)
}
