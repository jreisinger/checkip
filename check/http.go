package check

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const defaultHTTPTimeout = 10 * time.Second

// defaultHttpClient is reused by checks that make HTTP requests.
var defaultHttpClient = newHttpClient(&http.Client{Timeout: defaultHTTPTimeout})

type httpClient struct {
	client *http.Client
}

func newHttpClient(client *http.Client) httpClient {
	return httpClient{client: client}
}

func (c httpClient) Get(apiUrl string, headers map[string]string, queryParams map[string]string) ([]byte, error) {
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
		return nil, &httpStatusError{StatusCode: resp.StatusCode, url: apiUrl, status: resp.Status}
	}
	return body, nil
}

// httpStatusError is returned by Get when the response status is not 2xx, so
// callers can check the status code with errors.As instead of matching the
// error string.
type httpStatusError struct {
	StatusCode int
	url        string
	status     string
}

func (e *httpStatusError) Error() string {
	return fmt.Sprintf("GET %s: %s", e.url, e.status)
}

func (c httpClient) GetJson(apiUrl string, headers map[string]string, queryParams map[string]string, response interface{}) error {
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
