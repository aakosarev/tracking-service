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
	h.Logger.Info("Put tracking data in database")
	w.Header().Set("Content-Type", "application/json")
	trackingNumber := r.URL.Query().Get("tracking_number")
	courierCode := r.URL.Query().Get("courier_code")
	if trackingNumber == "" || courierCode == "" {
		w.WriteHeader(400)
		w.Write([]byte(`{"message":"Request type error
						            Please check the API documentation for the request type of this API"}`))
		return
	}
	inputData := InputData{TrackingNumber: trackingNumber, CourierCode: courierCode}
	buf, err := json.Marshal(inputData)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(`{"message":"Server error
									Please contact us: kosarevjob@gmail.com"}`))
		return
	}
	cfg := config.GetConfig()
	reqURL := fmt.Sprintf("%s%s", cfg.TrackingMore.BaseUrl, "/v3/trackings/create")
	req1, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(buf))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(`{"message":"Server error
									Please contact us: kosarevjob@gmail.com"}`))
		return
	}
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("Tracking-Api-Key", cfg.TrackingMore.ApiKey)
	defer req1.Body.Close()

	httpClient := &http.Client{}

	resp1, _ := httpClient.Do(req1)
	defer resp1.Body.Close()

	m := map[string]interface{}{}
	if err := json.NewDecoder(resp1.Body).Decode(&m); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(`{"message":"Server error
									Please contact us: kosarevjob@gmail.com"}`))
		return
	}

	status := int(m["code"].(float64))
	h.Logger.Debug(status)

	if status == http.StatusOK || status == 423 {
		reqURL = fmt.Sprintf("%s%s", cfg.TrackingMore.BaseUrl, fmt.Sprintf("/v3/trackings/get?tracking_numbers=%s", trackingNumber))
		req2, err := http.NewRequest(http.MethodGet, reqURL, nil)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(`{"message":"Server error
									    Please contact us: kosarevjob@gmail.com"}`))
			return
		}
		req2.Header.Set("Content-Type", "application/json")
		req2.Header.Set("Tracking-Api-Key", cfg.TrackingMore.ApiKey)
		defer req1.Body.Close()

		resp2, _ := httpClient.Do(req2)
		defer resp2.Body.Close()

		trackingResult := TrackingResult{}
		if err = json.NewDecoder(resp2.Body).Decode(&trackingResult); err != nil {
			http.Error(w, fmt.Sprintf("Server error 5: %s", err.Error()), 500)
			w.Write([]byte(`{"message":"Server error
									    Please contact us: kosarevjob@gmail.com"}`))
			return
		}

		status := int(trackingResult.Code)
		h.Logger.Debug(status)

		if status != http.StatusOK {
			w.WriteHeader(status)
			w.Write([]byte(fmt.Sprintf(`{"message":"Message from TrackingMore: %s"}`, trackingResult.Message)))
			return
		}

		body, _ := json.Marshal(trackingResult) // TODO: Delete this
		h.Logger.Debug(string(body))            // TODO: Delete this

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Request response is successful"}`))
		return
	} else {
		w.WriteHeader(status)
		w.Write([]byte(fmt.Sprintf(`{"message":"Message from TrackingMore: %s"}`, m["message"])))
		return
	}
}
