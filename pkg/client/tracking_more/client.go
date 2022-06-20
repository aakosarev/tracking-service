package tracking_more

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

func NewClient(httpClient *http.Client, baseURL string, apiKey string) *Client {
	return &Client{httpClient: httpClient, baseURL: baseURL, apiKey: apiKey}
}

func (c *Client) doRequest(ctx context.Context, method, endpoint string, in interface{}) (resp *http.Response, err error) {
	var jsonStr []byte
	if in != nil {
		jsonStr, err = json.Marshal(in)
	}
	if err != nil {
		return nil, err
	}
	requestUrl := fmt.Sprintf("%s%s", c.baseURL, endpoint)

	req, err := http.NewRequestWithContext(ctx, method, requestUrl, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Tracking-Api-Key", c.apiKey)
	defer req.Body.Close()
	return c.httpClient.Do(req)
}

func (c *Client) CreateTracking(ctx context.Context, inputData *InputDataForCreatingTracking) (resp *http.Response, err error) {
	resp, err = c.doRequest(ctx, http.MethodPost, "/v3/trackings/create", inputData)
	return
}

func (c *Client) GetResult(ctx context.Context, trackingNumber string) (resp *http.Response, err error) {
	resp, err = c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/v3/trackings/get?tracking_numbers=%s", trackingNumber), nil)
	return
}
