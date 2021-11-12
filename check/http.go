package check

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// DefaultHttpClient is reused by checks that make HTTP requests.
var DefaultHttpClient = HttpClient{client: &http.Client{Timeout: 5 * time.Second}}

type HttpClient struct {
	client *http.Client
}

func (c HttpClient) Get(apiUrl string, headers map[string]string, queryParams map[string]string) ([]byte, error) {
	apiURL, err := url.Parse(apiUrl)
	if err != nil {
		return nil, err
	}

	// Set query parameters.
	if len(queryParams) != 0 {
		vals := url.Values{}
		for k, v := range queryParams {
			vals.Add(k, v)
		}
		apiURL.RawQuery = vals.Encode()
	}

	req, err := http.NewRequest("GET", apiURL.String(), nil)
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

func (c HttpClient) GetJson(apiUrl string, headers map[string]string, queryParams map[string]string, response interface{}) error {
	b, err := c.Get(apiUrl, headers, queryParams)
	if err != nil {
		return err
	}
	if response != nil {
		if err := json.Unmarshal(b, response); err != nil {
			return fmt.Errorf("unmarshalling JSON from %s: %v", apiUrl, err)
		}
	}
	return nil
}
