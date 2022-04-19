package check

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/jreisinger/checkip"
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
		na(m.Country), na(m.IsoCode), na(m.City), m.IsInEU)
}

func (m maxmind) JsonString() (string, error) {
	b, err := json.Marshal(m)
	return string(b), err
}

// MaxMind gets geolocation data from maxmind.com's GeoLite2-City.mmdb.
func MaxMind(ip net.IP) (checkip.Result, error) {
	result := checkip.Result{
		Name: "maxmind.com",
		Type: checkip.TypeInfo,
	}

	licenseKey, err := getConfigValue("MAXMIND_LICENSE_KEY")
	if err != nil {
		return result, newCheckError(err)
	}
	if licenseKey == "" {
		return result, nil
	}

	// file := "/var/tmp/GeoLite2-City.mmdb"
	file, err := getDbFilesPath("GeoLite2-City.mmdb")
	if err != nil {
		return result, err
	}

	url := "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=" + licenseKey + "&suffix=tar.gz"

	if err := updateFile(file, url, "tgz"); err != nil {
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

	result.Info = maxmind{
		City:    geo.City.Names["en"],
		Country: geo.Country.Names["en"],
		IsoCode: geo.Country.IsoCode,
		IsInEU:  geo.Country.IsInEuropeanUnion,
	}

	return result, nil
}
