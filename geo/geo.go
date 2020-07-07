package geo

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jreisinger/checkip/util"
	"github.com/oschwald/geoip2-golang"
)

// DB represents MaxMind's GeoIP database.
type DB struct {
	Filepath string
	URL      string
	DB       *geoip2.Reader
	Location []string
}

// New creates GeoDB with some defaults.
func New() *DB {
	return &DB{
		Filepath: "/var/tmp/GeoLite2-City.mmdb",
	}
}

// Update downloads and creates database file if not present,
// updates if file is older than a week.
func (g *DB) Update() error {
	licenseKey := os.Getenv("GEOIP_LICENSE_KEY")
	g.URL = "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=" + licenseKey + "&suffix=tar.gz"
	file, err := os.Stat(g.Filepath)

	if os.IsNotExist(err) {
		if licenseKey == "" {
			return fmt.Errorf("environment variable GEOIP_LICENSE_KEY is not set")
		}

		r, err := util.DownloadFile(g.URL)
		if err != nil {
			return err
		}
		if err := util.ExtractFile(g.Filepath, r); err != nil {
			return err
		}

		return nil // don't check ModTime if file does not exist
	}

	if util.IsOlderThanOneWeek(file.ModTime()) {
		if licenseKey == "" {
			log.Printf("environment variable GEOIP_LICENSE_KEY is not set")
			return nil
		}

		r, err := util.DownloadFile(g.URL)
		if err != nil {
			return err
		}
		if err := util.ExtractFile(g.Filepath, r); err != nil {
			return err
		}
	}

	return nil
}

// Open loads database from file to memory.
func (g *DB) Open() error {
	db, err := geoip2.Open(g.Filepath)
	if err != nil {
		return err
	}
	g.DB = db
	return nil
}

// Close closes database file.
func (g *DB) Close() {
	g.DB.Close()
}

// ForIP fills the geolocation data into the GeoDB struct.
func (g *DB) ForIP(ip net.IP) error {
	if err := g.Update(); err != nil {
		return fmt.Errorf("can't update geo DB: %v", err)
	}

	if err := g.Open(); err != nil {
		return fmt.Errorf("can't load geo DB: %v", err)
	}
	defer g.Close()

	record, err := g.DB.City(ip)
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
