package check

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/TylerBrock/colorjson"
	"io"
	"net/http"
	"net/url"
)

/*
to be moved to http.go
*/

var mispURL = ""

func (c httpClient) Post(apiUrl string, headers map[string]string, query map[string]string) ([]byte, error) {
	apiURL, err := url.Parse(apiUrl)
	if err != nil {
		return nil, err
	}

	jsonByte, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", apiURL.String(), bytes.NewBuffer(jsonByte))
	if err != nil {
		return nil, err
	}

	// Set request headers.
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("GET %s: %s", apiUrl, resp.Status)
	}
	return body, nil
}

func (c httpClient) PostJson(apiUrl string, headers map[string]string, query map[string]string, response interface{}) error {
	b, err := c.Post(apiUrl, headers, query)
	if err != nil {
		return err
	}

	if Debug {
		var dat map[string]interface{}
		json.Unmarshal(b, &dat)
		// Make a custom formatter with indent set
		f := colorjson.NewFormatter()
		f.Indent = 4
		// Marshall the Colorized JSON
		s, _ := f.Marshal(dat)
		fmt.Println(string(s))
	}
	if response != nil {
		if err := json.Unmarshal(b, response); err != nil {
			return fmt.Errorf("unmarshalling JSON from %s: %v", apiUrl, err)
		}
	}
	return nil
}

/*
===================================
*/

type misp struct {
	Response response `json:"response"`
}

type response struct {
	Attribute attribute `json:"attribute"`
}

type attribute []struct {
	ID      string `json:"id"`
	EventID string `json:"event_id"`
	Comment string `json:"comment"`
	Deleted bool   `json:"deleted"`
	Value   string `json:"value"`
	Event   event  `json:"Event"`
}

type event struct {
	PublishTimestamp string `json:"publish_timestamp"`
	ID               string `json:"id"`
	OrgCID           string `json:"orgc_id"`
	Info             string `json:"info"`
	UUID             string `json:"uuid"`
}

// Misp gets generic information from search.misp.io.
func Misp(ipaddr net.IP) (Check, error) {

	result := Check{Description: "Misp", Type: InfoAndIsMalicious}

	headers := map[string]string{
		"Accept":       "application/vnd.misp.api.v3.host.v1+json",
		"Content-Type": "application/x-www-form-urlencoded;charset=UTF-8",
	}

	// mandatory MISP_KEY
	apiKey, err := getConfigValue("MISP_KEY")
	if err != nil {
		return result, newCheckError(err)
	}
	if apiKey == "" {
		result.MissingCredentials = "MISP_KEY"
		return result, nil
	}
	headers["Accept"] = "application/json"
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = apiKey

	// mandatory MISP_URL
	if mispURL == "" { // or set by test
		var err1 error
		mispURL, err1 = getConfigValue("MISP_URL")
		if err1 != nil {
			return result, newCheckError(err1)
		}
	}
	if mispURL == "" {
		result.MissingCredentials = "MISP_URL"
		return result, nil
	}

	// optional MISP_OPT
	options, _ := getConfigValue("MISP_OPT")
	var re1 = regexp.MustCompile(`selfsigned`)
	var re2 = regexp.MustCompile(`\d+[h|d]`)

	self := re1.MatchString(options)
	last := re2.FindString(options)

	jsonQuery := fmt.Sprintf(`{
		"returnFormat": "json", 
		"type": "ip-src", 
		"value": "%s", 
		"limit": 5
		}`, ipaddr)

	var dat map[string]string
	json.Unmarshal([]byte(jsonQuery), &dat)

	if last != "" {
		dat["last"] = last
	}

	var misp misp
	apiURL := fmt.Sprintf("%s/attributes/restSearch", mispURL)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	var httpClient = defaultHttpClient
	if self {
		httpClient = newHttpClient(&http.Client{Timeout: 5 * time.Second, Transport: tr})
	}

	if err := httpClient.PostJson(apiURL, headers, dat, &misp); err != nil {
		return result, newCheckError(err)
	}
	result.IpAddrInfo = misp
	result.IpAddrIsMalicious = true

	return result, nil
}

// Summary returns interesting information from the check.
func (c misp) Summary() string {
	cleanCountry := ` \([A-Z]{2}\)`
	var re = regexp.MustCompile(cleanCountry)
	cleanASN := ` AS\d+ - .*`
	var re2 = regexp.MustCompile(cleanASN)

	res := make(map[string][]string)
	for _, e := range c.Response.Attribute {
		edesc := fmt.Sprintf("* [%s orgc: %s] %s", e.Event.ID, e.Event.OrgCID, e.Event.Info)
		s := re.ReplaceAllString(e.Comment, "")
		s = re2.ReplaceAllString(s, "")
		res[edesc] = append(res[edesc], s)
	}

	sum := ""
	for key := range res {
		sum = sum + key + ": " + strings.Join(res[key], ", ")
	}
	return sum
}

func (c misp) Json() ([]byte, error) {
	return json.Marshal(c)
}
