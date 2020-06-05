package geodb

import (
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

func isOlderThanOneWeek(t time.Time) bool {
	return time.Now().Sub(t) > 7*24*time.Hour
}

// Update downloads and creates database file if not present,
// updates if file is older than a week.
func (g *GeoDB) Update(url string) error {
	if file, err := os.Stat(g.Filepath); os.IsNotExist(err) || isOlderThanOneWeek(file.ModTime()) {
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

// GetLocation fills the geolocation data into the GeoDB struct.
func (g *GeoDB) GetLocation(ip net.IP) error {
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
