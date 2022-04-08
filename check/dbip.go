package check

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/jreisinger/checkip"
	"github.com/oschwald/geoip2-golang"
)

type dbip struct {
	City    string `json:"city"`
	Country string `json:"country"`
	IsoCode string `json:"iso_code"`
	IsInEU  bool   `json:"is_in_eu"`
}

func (d dbip) Summary() string {
	return fmt.Sprintf("country: %s (%s), city: %s, EU member: %t",
		na(d.Country), na(d.IsoCode), na(d.City), d.IsInEU)
}

func (d dbip) JsonString() (string, error) {
	b, err := json.Marshal(d)
	return string(b), err
}

// DBip gets geolocation data from https://db-ip.com/db/download/ip-to-city-lite
func DBip(ip net.IP) (checkip.Result, error) {
	result := checkip.Result{
		Name: "db-ip.com",
		Type: checkip.TypeInfo,
	}

	file := "/var/tmp/dbip-city-lite.mmdb"
	url := fmt.Sprintf(
		"https://download.db-ip.com/free/dbip-city-lite-%s.mmdb.gz",
		time.Now().Format("2006-01"),
	)

	if err := updateFile(file, url, "gz"); err != nil {
		return result, newCheckError(err)
	}

	db, err := geoip2.Open(file)
	if err != nil {
		return result, newCheckError(fmt.Errorf("can't load DB file: %v", err))
	}
	defer db.Close()

	geo, err := db.City(ip)
	if err != nil {
		return result, newCheckError(err)
	}

	result.Info = dbip{
		City:    geo.City.Names["en"],
		Country: geo.Country.Names["en"],
		IsoCode: geo.Country.IsoCode,
		IsInEU:  geo.Country.IsInEuropeanUnion,
	}

	return result, nil
}
