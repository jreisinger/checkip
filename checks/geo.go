package checks

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/jreisinger/checkip/check"
	"github.com/oschwald/geoip2-golang"
)

type geo struct {
	City    string `json:"city"`
	Country string `json:"country"`
	IsoCode string `json:"iso_code"`
}

func (g geo) Summary() string {
	return fmt.Sprintf("city: %s, country: %s (%s)", check.Na(g.City), check.Na(g.Country), check.Na(g.IsoCode))
}

func (g geo) JsonString() (string, error) {
	b, err := json.Marshal(g)
	return string(b), err
}

// CheckGeo gets geolocation data from maxmind.com's GeoLite2-City.mmdb.
func CheckGeo(ip net.IP) (check.Result, error) {
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

	record, err := db.City(ip)
	if err != nil {
		return check.Result{}, check.NewError(err)
	}

	geo := geo{
		City:    record.City.Names["en"],
		Country: record.Country.Names["en"],
		IsoCode: record.Country.IsoCode,
	}

	return check.Result{
		Name: "maxmind.com",
		Type: check.TypeInfo,
		Info: geo,
	}, nil
}
