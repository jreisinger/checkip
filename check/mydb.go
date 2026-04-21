package check

import (
	"encoding/json"
	"fmt"
	"net"
)

type myDB struct {
	Alert     bool   `json:"alert"`
	Registred bool   `json:"registred"`
	Created   string `json:"created"`
	Excluded  bool   `json:"excluded"`
	Survey    string `json:"survey"`
}

var myDBUrl = ""

// MyDB gets generic information from api.myDB.io.
func MyDB(ipaddr net.IP) (Check, error) {

	result := Check{Description: "MyDB", Type: InfoAndIsMalicious}

	if myDBUrl == "" { // or set by test
		var err1 error
		myDBUrl, err1 = getConfigValue("MYDB_URL")
		if err1 != nil {
			return result, newCheckError(err1)
		}
	}
	if myDBUrl == "" {
		result.MissingCredentials = "MYDB_URL"
		return result, nil
	}

	apiKey, err := getConfigValue("MYDB_API_KEY")
	if err != nil {
		return result, newCheckError(err)
	}
	if apiKey == "" {
		result.MissingCredentials = "MYDB_API_KEY"
		return result, nil
	}

	headers := map[string]string{
		"Authorization": "bearer " + apiKey,
		"Accept":        "application/json",
		"Content-Type":  "application/x-www-form-urlencoded",
	}
	var t myDB
	apiURL := fmt.Sprintf("%s/%s", myDBUrl, ipaddr)
	if err := defaultHttpClient.GetJson(apiURL, headers, map[string]string{}, &t); err != nil {
		return result, newCheckError(err)
	}

	result.IpAddrIsMalicious = t.Alert

	result.IpAddrInfo = t

	return result, nil
}

func (o myDB) Summary() string {
	sum := ""
	if o.Registred {
		sum = fmt.Sprintf("Registred on %s", o.Created)
	}
	if o.Excluded {
		sum = fmt.Sprintf(" [excluded]")
	}
	if o.Survey != "" {
		sum = fmt.Sprintf("%s ** Survey: %s", sum, o.Survey)
	}
	return sum

}

func (o myDB) Json() ([]byte, error) {
	return json.Marshal(o)
}
