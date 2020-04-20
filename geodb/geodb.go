package geodb

import (
	"os"

	"github.com/oschwald/geoip2-golang"
)

// GeoDB represents MaxMind's GeoIP database.
type GeoDB struct {
	Filepath string
	URL		 string
	Age      string
	DB       *geoip2.Reader
}

// Update downloads and creates database file if not present,
// updates if file is older than a week.
func (g *GeoDB) Update(url string) error {
	if _, err := os.Stat(g.Filepath); os.IsNotExist(err) {
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

// Load loads database from file to memory.
func (g *GeoDB) Load() error {
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