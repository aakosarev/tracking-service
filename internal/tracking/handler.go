package tracking

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aakosarev/tracking-service/internal/config"
	"github.com/aakosarev/tracking-service/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Handler struct {
	Logger  *logging.Logger
	Service Service
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, "/tracking", h.Track)
}

func (h *Handler) Track(w http.ResponseWriter, r *http.Request) {
	trackingNumber := r.URL.Query().Get("tracking_number")
	courierCode := r.URL.Query().Get("courier_code")
	if trackingNumber == "" || courierCode == "" {
		http.Error(w, "Request type error", 400)
		return
	}
	inputData := InputData{TrackingNumber: trackingNumber, CourierCode: courierCode}
	buf, err := json.Marshal(inputData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Server error: %s", err.Error()), 511)
		return
	}
	cfg := config.GetConfig()
	reqURL := fmt.Sprintf("%s%s", cfg.TrackingMore.BaseUrl, "/v3/trackings/create")
	req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(buf))
	if err != nil {
		http.Error(w, fmt.Sprintf("Server error: %s", err.Error()), 511)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Tracking-Api-Key", cfg.TrackingMore.ApiKey)
	defer req.Body.Close()

	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK || resp.StatusCode == 423 {
		reqURL = fmt.Sprintf("%s%s", cfg.TrackingMore.BaseUrl, fmt.Sprintf("/v3/trackings/get?tracking_numbers=%s", trackingNumber))
		req, err = http.NewRequest(http.MethodGet, reqURL, nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("Server error: %s", err.Error()), 511)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Tracking-Api-Key", cfg.TrackingMore.ApiKey)
		resp, err = httpClient.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		trackingResult := TrackingResult{}
		if err = json.NewDecoder(resp.Body).Decode(&trackingResult); err != nil {
			fmt.Println(err)
		}
		body, err := json.Marshal(trackingResult)
		if err != nil {
			fmt.Println(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(body)

	} else {
		w.WriteHeader(resp.StatusCode)
	}
}
