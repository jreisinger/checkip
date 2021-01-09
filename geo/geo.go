package geo

import (
	"fmt"
	"net"

	"github.com/jreisinger/checkip/util"
	"github.com/oschwald/geoip2-golang"
)

// DB represents MaxMind's GeoIP database.
type DB struct {
	Location []string
}

// New creates GeoDB with some defaults.
func New() *DB {
	return &DB{}
}

// ForIP fills the geolocation data into the GeoDB struct.
func (g *DB) ForIP(ip net.IP) error {
	licenseKey, err := util.GetConfigValue("GEOIP_LICENSE_KEY")
	if err != nil {
		return fmt.Errorf("getting licence key: %w", err)
	}

	file := "/var/tmp/GeoLite2-City.mmdb"
	url := "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=" + licenseKey + "&suffix=tar.gz"

	if err := util.Update(file, url); err != nil {
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

	return nil
}
