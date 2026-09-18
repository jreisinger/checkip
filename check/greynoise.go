package check

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
)

var greynoiseurl = "https://api.greynoise.io/v3/community/"

type grey struct {
	IP             string `json:"ip"`
	Noise          bool   `json:"noise"`
	Riot           bool   `json:"riot"`
	Classification string `json:"classification"`
	Name           string `json:"name"`
	Link           string `json:"link"`
	LastSeen       string `json:"last_seen"`
	Message        string `json:"message"`
}

// Json implements IpInfo.
func (g grey) Json() ([]byte, error) {
	return json.Marshal(g)
}

// Summary implements IpInfo.
func (g grey) Summary() string {
	return na(g.Message)
}

// GreyNoise is a check for GreyNoise.
func GreyNoise(ipaddr net.IP) (Check, error) {
	result := Check{
		Description: "greynoise.io",
		Type:        InfoAndIsMalicious,
	}

	var response grey

	apiURL := greynoiseurl + ipaddr.String()
	headers := map[string]string{"accept": "application/json"}

	if err := defaultHttpClient.GetJson(apiURL, headers, map[string]string{}, &response); err != nil {
		var statusErr *httpStatusError
		if errors.As(err, &statusErr) && statusErr.StatusCode == http.StatusNotFound {
			result.IpAddrInfo = grey{
				IP:      ipaddr.String(),
				Message: "IP not observed scanning the internet or contained in RIOT data set.",
			}
			return result, nil
		}
		return result, newCheckError(err)
	}

	result.IpAddrIsMalicious = response.Classification == "malicious"
	result.IpAddrInfo = response

	return result, nil
}
