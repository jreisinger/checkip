package check

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/oschwald/geoip2-golang"
)

type dbip struct {
	City    string `json:"city"`
	Country string `json:"country"`
	IsoCode string `json:"iso_code"`
	IsInEU  bool   `json:"is_in_eu"`
}

func (d dbip) Summary() string {
	// Get just non-empty strings.
	var parts []string
	for _, s := range []string{d.City, d.Country} {
		if strings.TrimSpace(s) != "" {
			parts = append(parts, s)
		}
	}

	return strings.Join(parts, ", ")
}

func (d dbip) Json() ([]byte, error) {
	return json.Marshal(d)
}

// DBip gets geolocation from https://db-ip.com/db/download/ip-to-city-lite.
func DBip(ip net.IP) (Check, error) {
	result := Check{
		Description: "db-ip.com",
		Type:        Info,
	}

	// file := "/var/tmp/dbip-city-lite.mmdb"
	file, err := getCachePath("dbip-city-lite.mmdb")
	if err != nil {
		return result, err
	}

	url := fmt.Sprintf(
		"https://download.db-ip.com/free/dbip-city-lite-%s.mmdb.gz",
		time.Now().Format("2006-01"),
	)

	if err := updateFile(file, url, "gz"); err != nil {
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

	result.IpAddrInfo = dbip{
		City:    geo.City.Names["en"],
		Country: geo.Country.Names["en"],
		IsoCode: geo.Country.IsoCode,
		IsInEU:  geo.Country.IsInEuropeanUnion,
	}

	return result, nil
}
