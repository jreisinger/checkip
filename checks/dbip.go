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
	Country   string `json:"country"`
	IsoCode   string `json:"iso_code"`
	Continent string `json:"continent"`
	IsInEU    bool   `json:"is_in_eu"`
}

func (d dbip) Summary() string {
	return fmt.Sprintf("country: %s (%s), continent: %s, EU member: %t",
		check.Na(d.Country), check.Na(d.IsoCode), d.Continent, d.IsInEU)
}

func (d dbip) JsonString() (string, error) {
	b, err := json.Marshal(d)
	return string(b), err
}

func DBip(ip net.IP) (check.Result, error) {
	file := "/var/tmp/dbip-country-lite.mmdb"
	url := fmt.Sprintf(
		"https://download.db-ip.com/free/dbip-country-lite-%s.mmdb.gz",
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

	record, err := db.Country(ip)
	if err != nil {
		return check.Result{}, check.NewError(err)
	}

	d := dbip{
		Country:   record.Country.Names["en"],
		IsoCode:   record.Country.IsoCode,
		Continent: record.Continent.Names["en"],
		IsInEU:    record.Country.IsInEuropeanUnion,
	}

	return check.Result{
		Name: "db-ip.com",
		Type: check.TypeInfo,
		Info: d,
	}, nil
}
