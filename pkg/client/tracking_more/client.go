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

func (c *Client) doRequest(ctx context.Context, method, endpoint string, in, out interface{}) (err error) {
	jsonStr, err := json.Marshal(in)
	if err != nil {
		return err
	}
	requestUrl := fmt.Sprintf("%s%s", c.baseURL, endpoint)

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

type InputData struct {
	TrackingNumber string `json:"tracking_number"`
	CourierCode    string `json:"courier_code"`
}

func (c *Client) CreateTracker(inputData *InputData) (out interface{}, err error) {
	err = c.doRequest(context.Background(), http.MethodPost, "/v3/trackings/create", inputData, &out)
	return
}
