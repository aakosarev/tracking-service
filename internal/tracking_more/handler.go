package tracking_more

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aakosarev/tracking-service/internal/config"
	"github.com/aakosarev/tracking-service/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

type Service interface {
	CreateOrUpdate(databaseData DatabaseData) error
	UpdateAll() ([]InputData, error)
}

type handler struct {
	service Service
	logger  *logging.Logger
}

func NewHandler(service Service, logger *logging.Logger) *handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, "/tracking_more", h.Track)
	router.HandlerFunc(http.MethodPost, "/update", h.UpdateAll)
}

func (h *handler) Track(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Put tracking_more data in database")
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
		h.logger.Debug(err)
		w.WriteHeader(500)
		w.Write([]byte(`{"message":"Server error. Please contact us: kosarevjob@gmail.com"}`))
		return
	}
	cfg := config.GetConfig()
	reqURL := "https://api.trackingmore.com/v3/trackings/realtime"
	req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(buf))
	if err != nil {
		h.logger.Debug(err)
		w.WriteHeader(500)
		w.Write([]byte(`{"message":"Server error. Please contact us: kosarevjob@gmail.com"}`))
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Tracking-Api-Key", cfg.TrackingMore.ApiKey)
	defer req.Body.Close()

	httpClient := &http.Client{}

	resp, _ := httpClient.Do(req)
	if err != nil {
		h.logger.Debug(err)
		w.WriteHeader(500)
		w.Write([]byte(`{"message":"Server error. Please contact us: kosarevjob@gmail.com"}`))
		return
	}
	defer resp.Body.Close()

	trackingResult := TrackingResult{}
	if err = json.NewDecoder(resp.Body).Decode(&trackingResult); err != nil {
		h.logger.Debug(err)
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
	if err = h.service.CreateOrUpdate(databaseData); err != nil {
		h.logger.Debug(err)
		w.WriteHeader(500)
		w.Write([]byte(`{"message":"Server error. Please contact us: kosarevjob@gmail.com"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Request response is successful"}`))
	return
}

func (h *handler) UpdateAll(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Update all items into database")
	w.Header().Set("Content-Type", "application/json")
	inputData, err := h.service.UpdateAll()
	if err != nil {
		h.logger.Debug(err)
		w.WriteHeader(500)
		w.Write([]byte(`{"message":"Server error. Please contact us: kosarevjob@gmail.com"}`))
		return
	}

	httpClient := &http.Client{}

	for _, v := range inputData {
		buf, _ := json.Marshal(v)
		cfg := config.GetConfig()
		reqURL := "https://api.trackingmore.com/v3/trackings/realtime"
		req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(buf))
		if err != nil {
			h.logger.Debug(err)
			w.WriteHeader(500)
			w.Write([]byte(`{"message":"Server error. Please contact us: kosarevjob@gmail.com"}`))
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Tracking-Api-Key", cfg.TrackingMore.ApiKey)

		resp, _ := httpClient.Do(req)
		if err != nil {
			h.logger.Debug(err)
			w.WriteHeader(500)
			w.Write([]byte(`{"message":"Server error. Please contact us: kosarevjob@gmail.com"}`))
			return
		}

		trackingResult := TrackingResult{}
		if err = json.NewDecoder(resp.Body).Decode(&trackingResult); err != nil {
			h.logger.Debug(err)
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
		if err = h.service.CreateOrUpdate(databaseData); err != nil {
			h.logger.Debug(err)
			w.WriteHeader(500)
			w.Write([]byte(`{"message":"Server error. Please contact us: kosarevjob@gmail.com"}`))
			return
		}
		time.Sleep(1 * time.Second) // in order not to exceed the limit TrackingMore
		h.logger.Debug("Track number :", v.TrackingNumber, "successfully updated")
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Request response is successful"}`))
	return
}
