package tracking_more

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	apiKey     string
}

func NewClient(httpClient *http.Client, baseURL *url.URL, apiKey string) *Client {
	return &Client{httpClient: httpClient, baseURL: baseURL, apiKey: apiKey}
}

func (c *Client) doRequest(ctx context.Context, method, endpoint string, in, out interface{}) error {
	jsonStr, err := json.Marshal(in)
	if err != nil {
		return err
	}
	requestUrl := fmt.Sprintf("%s%s", c.baseURL.String(), endpoint)

	req, err := http.NewRequestWithContext(ctx, method, requestUrl, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Tracking-Api-Key", c.apiKey)
	defer req.Body.Close()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(out); err != nil {
		return err
	}

	return nil
}
