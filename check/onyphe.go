package check

import (
	"encoding/json"
	"fmt"
	"net"
	"sort"
	"strings"
	"strconv"

	"github.com/jreisinger/checkip"
)

type onyphe struct {
	Text    string     `json:"text"`
	Vulns   []string   `json:"vulns"`
	Results onypheData `json:"results"`
}

type onypheData []struct {
	OS       string `json:"os"`
	OsVendor string `json:"osvendor"`
	//Version   string `json:"version"`
	Port      interface{} `json:"port"`
	Protocol  string      `json:"protocol"`
	Product   string      `json:"product"`
	Transport string      `json:"transport"` // tcp, udp
}

var onypheUrl = "https://www.onyphe.io/api/v2"

// Onyphe gets generic information from api.onyphe.io.
func Onyphe(ipaddr net.IP) (checkip.Result, error) {
	result := checkip.Result{
		Name: "onyphe.io",
		Type: checkip.TypeInfoSec,
	}

	apiKey, err := getConfigValue("ONYPHE_API_KEY")
	if err != nil {
		return result, newCheckError(err)
	}
	if apiKey == "" {
		return result, nil
	}

	headers := map[string]string{
		"Authorization": "bearer " + apiKey,
		"Accept":        "application/json",
		//"Content-Type": "application/x-www-form-urlencoded",
	}
	var onyphe onyphe
	apiURL := fmt.Sprintf("%s/simple/datascan/%s", onypheUrl, ipaddr)
	if err := defaultHttpClient.GetJson(apiURL, headers, map[string]string{}, &onyphe); err != nil {
		return result, newCheckError(err)
	}

	for _, d := range onyphe.Results {
		var port string
		switch v := d.Port.(type) {
		case string:
			port = v
		default:
			port = fmt.Sprintf("%d", v)
		}
		if port != "80" && port != "443" && port != "53" { // undecidable ports
			result.Malicious = true
		}
	}

	result.Info = onyphe

	return result, nil
}

type byPortO onypheData

func (x byPortO) Len() int { return len(x) }
func (x byPortO) Less(i, j int) bool {
	var portI int
	var portJ int
	switch v := x[i].Port.(type) {
	case string:
		portI, _ = strconv.Atoi(v)
    case float64:
	    portI = int(v)
    case int:
	    portI = v
	}
	return portI > portJ
}

func (x byPortO) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

// Info returns interesting information from the check.
func (o onyphe) Summary() string {
	var portInfo []string
	service := make(map[string]int)
	sort.Sort(byPortO(o.Results))
	for _, d := range o.Results {
		var os string
		if d.OS != "" {
			os = d.OS
		}

		var osvendor string
		if d.OsVendor != "" {
			osvendor = d.OsVendor
		}

		var product string
		if d.Product != "" {
			product = d.Product + " "
		}

		var port string
		switch v := d.Port.(type) {
		case float64:
			port = fmt.Sprintf("%d", int(v))
		case string:
			port = v
		}

		sport := fmt.Sprintf("%s/%s", d.Transport, port)
		service[sport]++

		if service[sport] > 1 {
			continue
		}

		if os == "" && osvendor == "" {
			portInfo = append(portInfo, fmt.Sprintf("%s %s/%s", d.Protocol, d.Transport, port))
		} else {
			ss := nonEmpty(os, osvendor)
			portInfo = append(portInfo, fmt.Sprintf("%s %s/%s %s(%s)", d.Protocol, d.Transport, port, product, strings.Join(ss, ", ")))
		}
	}

	return fmt.Sprintf("Open: %s", strings.Join(portInfo, ", "))
}

func (o onyphe) Json() ([]byte, error) {
	return json.Marshal(o)
}
