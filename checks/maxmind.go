package checks

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/jreisinger/checkip/check"
	"github.com/oschwald/geoip2-golang"
)

type maxmind struct {
	City    string `json:"city"`
	Country string `json:"country"`
	IsoCode string `json:"iso_code"`
	IsInEU  bool   `json:"is_in_eu"`
}

func (m maxmind) Summary() string {
	return fmt.Sprintf("country: %s (%s), city: %s, EU member: %t",
		check.Na(m.Country), check.Na(m.IsoCode), check.Na(m.City), m.IsInEU)
}

func (m maxmind) JsonString() (string, error) {
	b, err := json.Marshal(m)
	return string(b), err
}

// MaxMind gets geolocation data from maxmind.com's GeoLite2-City.mmdb.
func MaxMind(ip net.IP) (check.Result, error) {
	licenseKey, err := check.GetConfigValue("MAXMIND_LICENSE_KEY")
	if err != nil {
		return check.Result{}, check.NewError(err)
	}

	file := "/var/tmp/GeoLite2-City.mmdb"
	url := "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=" + licenseKey + "&suffix=tar.gz"

	if err := check.UpdateFile(file, url, "tgz"); err != nil {
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

	m := maxmind{
		City:    geo.City.Names["en"],
		Country: geo.Country.Names["en"],
		IsoCode: geo.Country.IsoCode,
		IsInEU:  geo.Country.IsInEuropeanUnion,
	}

	return check.Result{
		Name: "maxmind.com",
		Type: check.TypeInfo,
		Info: m,
	}, nil
}
