package geodb

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/oschwald/geoip2-golang"
)

// GeoDB represents MaxMind's GeoIP database.
type GeoDB struct {
	Filepath string
	URL      string
	DB       *geoip2.Reader
	Location []string
}

// New creates GeoDB with some defaults.
func New() *GeoDB {
	return &GeoDB{
		Filepath: "/var/tmp/GeoLite2-City.mmdb",
	}
}

func isOlderThanOneWeek(t time.Time) bool {
	return time.Now().Sub(t) > 7*24*time.Hour
}

// Update downloads and creates database file if not present,
// updates if file is older than a week.
func (g *GeoDB) Update() error {
	licenseKey := os.Getenv("GEOIP_LICENSE_KEY")
	g.URL = "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=" + licenseKey + "&suffix=tar.gz"
	file, err := os.Stat(g.Filepath)

	if os.IsNotExist(err) {
		if licenseKey == "" {
			return fmt.Errorf("%s is not present and environment variable GEOIP_LICENSE_KEY is not set", g.Filepath)
		}

		r, err := downloadFile(g.URL)
		if err != nil {
			return err
		}
		if err := extractFile(g.Filepath, r); err != nil {
			return err
		}

		return nil // don't check ModTime if file does not exist
	}

	if isOlderThanOneWeek(file.ModTime()) {
		if licenseKey == "" {
			log.Printf("warning %s is outdated and environment variable GEOIP_LICENSE_KEY is not set", g.Filepath)
			return nil
		}

		r, err := downloadFile(g.URL)
		if err != nil {
			return err
		}
		if err := extractFile(g.Filepath, r); err != nil {
			return err
		}
	}

	return nil
}

// Open loads database from file to memory.
func (g *GeoDB) Open() error {
	db, err := geoip2.Open(g.Filepath)
	if err != nil {
		return err
	}
	g.DB = db
	return nil
}

// Close closes database file.
func (g *GeoDB) Close() {
	g.DB.Close()
}

// ForIP fills the geolocation data into the GeoDB struct.
func (g *GeoDB) ForIP(ip net.IP) error {
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
