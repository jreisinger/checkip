package check

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"github.com/logrusorgru/aurora"
	"github.com/oschwald/geoip2-golang"
)

type loc struct {
	IP      string `json:"ip"`
	City    string `json:"city"`
	Country string `json:"country"`
	IsoCode string `json:"iso_code"`
	IsInEU  bool   `json:"is_in_eu"`
	ASN     string `json:"asn"`
}

func (l loc) Summary() string {
	au := aurora.NewAurora(true)
	flag, _ := emocode(l.IsoCode)
	return fmt.Sprintf("%s (%s)%s %s", l.IP, au.Green(l.IsoCode), flag, l.ASN)
}

func emocode(x string) (string, error) {
	if len(x) != 2 {
		return "", errors.New("country code must be two letters")
	}
	if x[0] < 'A' || x[0] > 'Z' || x[1] < 'A' || x[1] > 'Z' {
		return "", errors.New("invalid country code")
	}
	return string(0x1F1E6+rune(x[0])-'A') + string(0x1F1E6+rune(x[1])-'A'), nil
}

func (l loc) Json() ([]byte, error) {
	return json.Marshal(l)
}

// IOCLoc gets geolocation data from maxmind.com's GeoLite2-City.mmdb.
func IOCLoc(ip net.IP) (Check, error) {
	result := Check{
		Description: "IOCLoc",
		Type: Info,
	}

	licenseKey, err := getConfigValue("MAXMIND_LICENSE_KEY")
	if err != nil {
		return result, newCheckError(err)
	}
	if licenseKey == "" {
		result.MissingCredentials = "MAXMIND_LICENSE_KEY"
		return result, nil
	}

	var dbCity, dbASN *geoip2.Reader
	for _, f := range []string{"GeoLite2-City", "GeoLite2-ASN"} {

		// file := "/var/tmp/GeoLite2-City.mmdb"
		file, err := getCachePath(f + ".mmdb")
		if err != nil {
			return result, err
		}

		url := "https://download.maxmind.com/app/geoip_download?edition_id=" + f + "&license_key=" + licenseKey + "&suffix=tar.gz"

		if err := updateFile(file, url, "tgz"); err != nil {
			return result, newCheckError(err)
		}
		switch f {
		case "GeoLite2-City":

			dbCity, err = geoip2.Open(file)
			if err != nil {
				return result, newCheckError(fmt.Errorf("can't load DB file %s: %v", file, err))
			}
			defer dbCity.Close()
		case "GeoLite2-ASN":
			dbASN, err = geoip2.Open(file)
			if err != nil {
				return result, newCheckError(fmt.Errorf("can't load DB file %s: %v", file, err))
			}
			defer dbASN.Close()
		}
	}

	geo, err := dbCity.City(ip)
	if err != nil {
		return result, newCheckError(err)
	}
	geoAsn, err := dbASN.ASN(ip)
	if err != nil {
		return result, newCheckError(err)
	}

	result.IpAddrInfo = loc{
		IP:      ip.String(),
		City:    geo.City.Names["en"],
		Country: geo.Country.Names["en"],
		IsoCode: geo.Country.IsoCode,
		IsInEU:  geo.Country.IsInEuropeanUnion,
		ASN:     fmt.Sprintf(" AS%d - %s", geoAsn.AutonomousSystemNumber, geoAsn.AutonomousSystemOrganization),
	}

	return result, nil
}
