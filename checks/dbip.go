package checks

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/jreisinger/checkip/check"
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
		check.Na(d.Country), check.Na(d.IsoCode), check.Na(d.City), d.IsInEU)
}

func (d dbip) JsonString() (string, error) {
	b, err := json.Marshal(d)
	return string(b), err
}

// DBip gets geolocation data from https://db-ip.com/db/download/ip-to-city-lite
func DBip(ip net.IP) (check.Result, error) {
	file := "/var/tmp/dbip-city-lite.mmdb"
	url := fmt.Sprintf(
		"https://download.db-ip.com/free/dbip-city-lite-%s.mmdb.gz",
		time.Now().Format("2006-01"),
	)

	if err := check.UpdateFile(file, url, "gz"); err != nil {
		return check.Result{}, check.NewError(err)
	}

	db, err := geoip2.Open(file)
	if err != nil {
		return check.Result{}, check.NewError(fmt.Errorf("can't load DB file: %v", err))
	}
	defer db.Close()

	geo, err := db.City(ip)
	if err != nil {
		return check.Result{}, check.NewError(err)
	}

	d := dbip{
		City:    geo.City.Names["en"],
		Country: geo.Country.Names["en"],
		IsoCode: geo.Country.IsoCode,
		IsInEU:  geo.Country.IsInEuropeanUnion,
	}

	return check.Result{
		Name: "db-ip.com",
		Type: check.TypeInfo,
		Info: d,
	}, nil
}
