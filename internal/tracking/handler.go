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
		w.Write([]byte(`{"message":"Request type error. Please check the API documentation for the request type of this API"}`))
		return
	}
	inputData := InputData{TrackingNumber: trackingNumber, CourierCode: courierCode}
	buf, err := json.Marshal(inputData)
	if err != nil {
		w.WriteHeader(500)
		h.Logger.Debug("1", err)
		w.Write([]byte(`{"message":"Server error. Please contact us: kosarevjob@gmail.com"}`))
		return
	}
	cfg := config.GetConfig()
	reqURL := fmt.Sprintf("%s%s", cfg.TrackingMore.BaseUrl, "/v3/trackings/realtime")
	req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(buf))
	if err != nil {
		w.WriteHeader(500)
		h.Logger.Debug("2", err)
		w.Write([]byte(`{"message":"Server error. Please contact us: kosarevjob@gmail.com"}`))
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Tracking-Api-Key", cfg.TrackingMore.ApiKey)
	defer req.Body.Close()

	httpClient := &http.Client{}

	resp, _ := httpClient.Do(req)
	defer resp.Body.Close()

	trackingResult := TrackingResult{}
	if err = json.NewDecoder(resp.Body).Decode(&trackingResult); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(`{"message":"Server error. Please contact us: kosarevjob@gmail.com"}`))
		return
	}

	status := int(trackingResult.Code)
	if status != http.StatusOK {
		w.WriteHeader(status)
		w.Write([]byte(fmt.Sprintf(`{"message":"Message from TrackingMore: %s"}`, trackingResult.Message)))
		return
	}

	databaseData := trackingResult.ConvertToDatabaseData()
	if err = h.Service.CreateOrUpdate(databaseData); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(`{"message":"Server error. Please contact us: kosarevjob@gmail.com"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Request response is successful"}`))
	return
}
