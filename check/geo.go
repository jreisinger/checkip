package check

import (
	"fmt"
	"net"
	"strings"

	"github.com/jreisinger/checkip/util"
	"github.com/oschwald/geoip2-golang"
)

// Geo holds geographic position of an IP address from MaxMind's GeoIP database.
type Geo struct {
	Location []string
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

// Name returns the name of the check.
func (g *Geo) Name() string {
	return fmt.Sprint("Geo")
}

// String returns the result of the check.
func (g *Geo) String() string {
	return fmt.Sprintf("%s", strings.Join(g.Location, " | "))
}
