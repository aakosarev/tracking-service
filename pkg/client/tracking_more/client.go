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
	var jsonStr []byte
	if in != nil {
		jsonStr, err = json.Marshal(in)
	}
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

type InputDataForCreatingTracking struct {
	TrackingNumber string `json:"tracking_number"`
	CourierCode    string `json:"courier_code"`
}

func (c *Client) CreateTracking(inputData *InputDataForCreatingTracking) (out map[string]interface{}, err error) {
	err = c.doRequest(context.Background(), http.MethodPost, "/v3/trackings/create", inputData, &out)
	return
}

func (c *Client) GetResult(trackingNumber string) (out *TrackingResult, err error) {
	err = c.doRequest(context.Background(), http.MethodGet, fmt.Sprintf("/v3/trackings/get?tracking_numbers=%s", trackingNumber), nil, &out)
	return
}

type TrackingResult struct {
	Code int64 `json:"code"`
	Data []struct {
		ScheduledAddress string `json:"Scheduled_Address"`
		Archived         bool   `json:"archived"`
		Consignee        string `json:"consignee"`
		CourierCode      string `json:"courier_code"`
		CreatedAt        string `json:"created_at"`
		CustomerEmail    string `json:"customer_email"`
		CustomerName     string `json:"customer_name"`
		CustomerPhone    string `json:"customer_phone"`
		DeliveryStatus   string `json:"delivery_status"`
		Destination      string `json:"destination"`
		DestinationInfo  struct {
			ArrivedAbroadDate      string `json:"arrived_abroad_date"`
			ArrivedDestinationDate string `json:"arrived_destination_date"`
			CourierCode            string `json:"courier_code"`
			CourierPhone           string `json:"courier_phone"`
			CustomsReceivedDate    string `json:"customs_received_date"`
			DepartedAirportDate    string `json:"departed_airport_date"`
			DispatchedDate         string `json:"dispatched_date"`
			ReceivedDate           string `json:"received_date"`
			ReferenceNumber        string `json:"reference_number"`
			Trackinfo              []struct {
				CheckpointDate              string `json:"checkpoint_date"`
				CheckpointDeliveryStatus    string `json:"checkpoint_delivery_status"`
				CheckpointDeliverySubstatus string `json:"checkpoint_delivery_substatus"`
				Location                    string `json:"location"`
				TrackingDetail              string `json:"tracking_detail"`
			} `json:"trackinfo"`
			Weblink string `json:"weblink"`
		} `json:"destination_info"`
		DestinationTrackNumber string `json:"destination_track_number"`
		ExchangeNumber         string `json:"exchangeNumber"`
		ID                     string `json:"id"`
		LastestCheckpointTime  string `json:"lastest_checkpoint_time"`
		LatestEvent            string `json:"latest_event"`
		LogisticsChannel       string `json:"logistics_channel"`
		Note                   string `json:"note"`
		OrderNumber            string `json:"order_number"`
		OriginInfo             struct {
			ArrivedAbroadDate      string `json:"arrived_abroad_date"`
			ArrivedDestinationDate string `json:"arrived_destination_date"`
			CourierCode            string `json:"courier_code"`
			CourierPhone           string `json:"courier_phone"`
			CustomsReceivedDate    string `json:"customs_received_date"`
			DepartedAirportDate    string `json:"departed_airport_date"`
			DispatchedDate         string `json:"dispatched_date"`
			ReceivedDate           string `json:"received_date"`
			ReferenceNumber        string `json:"reference_number"`
			Trackinfo              []struct {
				CheckpointDate              string `json:"checkpoint_date"`
				CheckpointDeliveryStatus    string `json:"checkpoint_delivery_status"`
				CheckpointDeliverySubstatus string `json:"checkpoint_delivery_substatus"`
				Location                    string `json:"location"`
				TrackingDetail              string `json:"tracking_detail"`
			} `json:"trackinfo"`
			Weblink string `json:"weblink"`
		} `json:"origin_info"`
		Original              string `json:"original"`
		Previously            string `json:"previously"`
		ScheduledDeliveryDate string `json:"scheduled_delivery_date"`
		ServiceCode           string `json:"service_code"`
		ShippingDate          string `json:"shipping_date"`
		StatusInfo            string `json:"status_info"`
		StayTime              int64  `json:"stay_time"`
		Substatus             string `json:"substatus"`
		Title                 string `json:"title"`
		TrackingNumber        string `json:"tracking_number"`
		TransitTime           int64  `json:"transit_time"`
		UpdateDate            string `json:"update_date"`
		Updating              bool   `json:"updating"`
		Weight                string `json:"weight"`
	} `json:"data"`
	Message string `json:"message"`
}
